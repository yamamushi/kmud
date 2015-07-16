package session

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Cristofori/kmud/combat"
	"github.com/Cristofori/kmud/events"
	"github.com/Cristofori/kmud/model"
	"github.com/Cristofori/kmud/types"
	"github.com/Cristofori/kmud/utils"
	// "log"
	// "os"
	"strings"
	"time"
)

type Session struct {
	conn   io.ReadWriter
	user   types.User
	player types.PC
	room   types.Room

	prompt string
	states map[string]string

	userInputChannel chan string
	inputModeChannel chan userInputMode
	prompterChannel  chan utils.Prompter
	panicChannel     chan interface{}
	eventChannel     chan events.Event

	silentMode bool

	replyId types.Id

	// logger *log.Logger
}

func NewSession(conn io.ReadWriter, user types.User, player types.PC) *Session {
	var session Session
	session.conn = conn
	session.user = user
	session.player = player
	session.room = model.GetRoom(player.GetRoomId())

	session.prompt = "%h/%H> "
	session.states = map[string]string{}

	session.userInputChannel = make(chan string)
	session.inputModeChannel = make(chan userInputMode)
	session.prompterChannel = make(chan utils.Prompter)
	session.panicChannel = make(chan interface{})
	session.eventChannel = events.Register(player)

	session.silentMode = false

	// file, err := os.OpenFile(player.GetName()+".log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	// utils.PanicIfError(err)

	// session.logger = log.New(file, player.GetName()+" ", log.LstdFlags)

	model.Login(player)

	return &session
}

type userInputMode int

const (
	CleanUserInput userInputMode = iota
	RawUserInput   userInputMode = iota
)

func (self *Session) Exec() {
	defer events.Unregister(self.player)
	defer model.Logout(self.player)

	self.printLineColor(types.ColorWhite, "Welcome, "+self.player.GetName())
	self.printRoom()

	// Main routine in charge of actually reading input from the connection object,
	// also has built in throttling to limit how fast we are allowed to process
	// commands from the user.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				self.panicChannel <- r
			}
		}()

		throttler := utils.NewThrottler(200 * time.Millisecond)

		for {
			mode := <-self.inputModeChannel
			prompter := <-self.prompterChannel
			input := ""

			switch mode {
			case CleanUserInput:
				input = utils.GetUserInputP(self.conn, prompter, self.user.GetColorMode())
			case RawUserInput:
				input = utils.GetRawUserInputP(self.conn, prompter, self.user.GetColorMode())
			default:
				panic("Unhandled case in switch statement (userInputMode)")
			}

			throttler.Sync()
			self.userInputChannel <- input
		}
	}()

	// Main loop
	for {
		input := self.getUserInputP(RawUserInput, self)
		if input == "" || input == "logout" || input == "quit" {
			return
		}

		if strings.HasPrefix(input, "/") {
			self.handleCommand(utils.Argify(input[1:]))
		} else {
			self.handleAction(utils.Argify(input))
		}
	}
}

func (self *Session) printLineColor(color types.Color, line string, a ...interface{}) {
	self.user.WriteLine(types.Colorize(color, fmt.Sprintf(line, a...)))
}

func (self *Session) printLine(line string, a ...interface{}) {
	self.printLineColor(types.ColorWhite, line, a...)
}

func (self *Session) printError(err string, a ...interface{}) {
	self.printLineColor(types.ColorRed, err, a...)
}

func (self *Session) printRoom() {
	playerList := model.PlayerCharactersIn(self.room.GetId(), self.player.GetId())
	npcList := model.NpcsIn(self.room.GetId())
	area := model.GetArea(self.room.GetAreaId())

	self.printLine(self.room.ToString(playerList, npcList,
		model.GetItems(self.room.GetItems()), area))
}

func (self *Session) clearLine() {
	utils.ClearLine(self.conn)
}

func (self *Session) asyncMessage(message string) {
	self.clearLine()
	self.printLine(message)
}

// Same behavior as menu.Exec(), except that it uses getUserInput
// which doesn't block the event loop while waiting for input
func (self *Session) execMenu(menu *utils.Menu) (string, types.Id) {
	choice := ""
	var data types.Id

	for {
		menu.Print(self.conn, self.user.GetColorMode())
		choice = self.getUserInputP(CleanUserInput, menu)
		if menu.HasAction(choice) || choice == "" {
			data = menu.GetData(choice)
			break
		}

		if choice != "?" {
			self.printError("Invalid selection")
		}
	}
	return choice, data
}

// getUserInput allows us to retrieve user input in a way that doesn't block the
// event loop by using channels and a separate Go routine to grab
// either the next user input or the next event.
func (self *Session) getUserInputP(inputMode userInputMode, prompter utils.Prompter) string {
	self.inputModeChannel <- inputMode
	self.prompterChannel <- prompter

	for {
		select {
		case input := <-self.userInputChannel:
			return input
		case event := <-self.eventChannel:
			if self.silentMode {
				continue
			}

			switch e := event.(type) {
			case events.TellEvent:
				self.replyId = e.From.GetId()
			case events.TickEvent:
				if !combat.InCombat(self.player) {
					oldHps := self.player.GetHitPoints()
					self.player.Heal(5)
					newHps := self.player.GetHitPoints()

					if oldHps != newHps {
						self.clearLine()
						self.user.Write(prompter.GetPrompt())
					}
				}
			}

			message := event.ToString(self.player)
			if message != "" {
				self.asyncMessage(message)
				self.user.Write(prompter.GetPrompt())
			}

		case quitMessage := <-self.panicChannel:
			panic(quitMessage)
		}
	}
}

func (self *Session) getUserInput(inputMode userInputMode, prompt string) string {
	return self.getUserInputP(inputMode, utils.SimplePrompter(prompt))
}

func (self *Session) getRawUserInput(prompt string) string {
	return self.getUserInput(RawUserInput, prompt)
}

func (self *Session) GetPrompt() string {
	prompt := self.prompt
	prompt = strings.Replace(prompt, "%h", strconv.Itoa(self.player.GetHitPoints()), -1)
	prompt = strings.Replace(prompt, "%H", strconv.Itoa(self.player.GetHealth()), -1)

	if len(self.states) > 0 {
		states := make([]string, len(self.states))

		i := 0
		for key, value := range self.states {
			states[i] = fmt.Sprintf("%s:%s", key, value)
			i++
		}

		prompt = fmt.Sprintf("%s %s", states, prompt)
	}

	return types.Colorize(types.ColorWhite, prompt)
}

func (self *Session) currentZone() types.Zone {
	return model.GetZone(self.room.GetZoneId())
}

func (self *Session) handleAction(action string, args []string) {
	if len(args) == 0 {
		direction := types.StringToDirection(action)

		if direction != types.DirectionNone {
			if self.room.HasExit(direction) {
				newRoom, err := model.MoveCharacter(self.player, direction)
				if err == nil {
					self.room = newRoom
					self.printRoom()
				} else {
					self.printError(err.Error())
				}

			} else {
				self.printError("You can't go that way")
			}

			return
		}
	}

	handler, found := actions[action]

	if found {
		if handler.alias != "" {
			handler = actions[handler.alias]
		}
		handler.exec(self, args)
	} else {
		self.printError("You can't do that")
	}
}

func (self *Session) handleCommand(command string, args []string) {
	if command[0] == '/' && self.user.IsAdmin() {
		quickRoom(self, command[1:])
		return
	}

	handler, found := commands[command]

	if found {
		if handler.alias != "" {
			handler = commands[handler.alias]
		}
		handler.exec(self, args)
	} else {
		self.printError("Unrecognized command: %s", command)
	}
}
