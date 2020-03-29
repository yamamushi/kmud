package session

import (
	"fmt"
	"github.com/yamamushi/kmud-2020/color"
	"io"
	"sort"
	"strconv"

	"github.com/yamamushi/kmud-2020/combat"
	"github.com/yamamushi/kmud-2020/events"
	"github.com/yamamushi/kmud-2020/model"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
	// "log"
	// "os"
	"strings"
	"time"
)

type Session struct {
	conn io.ReadWriter
	user types.User
	pc   types.PC

	prompt string
	states map[string]string

	userInputChannel chan string
	inputModeChannel chan userInputMode
	prompterChannel  chan utils.Prompter
	panicChannel     chan interface{}
	eventChannel     chan events.Event

	silentMode bool
	replyId    types.Id
	lastInput  string

	// logger *log.Logger
}

func NewSession(conn io.ReadWriter, user types.User, pc types.PC) *Session {
	var session Session
	session.conn = conn
	session.user = user
	session.pc = pc

	session.prompt = "%h/%H> "
	session.states = map[string]string{}

	session.userInputChannel = make(chan string)
	session.inputModeChannel = make(chan userInputMode)
	session.prompterChannel = make(chan utils.Prompter)
	session.panicChannel = make(chan interface{})
	session.eventChannel = events.Register(pc)

	session.silentMode = false

	// file, err := os.OpenFile(pc.GetName()+".log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	// utils.PanicIfError(err)

	// session.logger = log.New(file, pc.GetName()+" ", log.LstdFlags)

	model.Login(pc)

	return &session
}

type userInputMode int

const (
	CleanUserInput userInputMode = iota
	RawUserInput   userInputMode = iota
)

func (s *Session) Exec() {
	defer events.Unregister(s.pc)
	defer model.Logout(s.pc)

	s.WriteLine("Welcome, " + s.pc.GetName())
	s.PrintRoom()

	// Main routine in charge of actually reading input from the connection object,
	// also has built in throttling to limit how fast we are allowed to process
	// commands from the user.
	go func() {
		defer func() {
			s.panicChannel <- recover()
		}()

		throttler := utils.NewThrottler(200 * time.Millisecond)

		for {
			mode := <-s.inputModeChannel
			prompter := <-s.prompterChannel
			input := ""

			switch mode {
			case CleanUserInput:
				input = utils.GetUserInputP(s.conn, prompter, s.user.GetColorMode())
			case RawUserInput:
				input = utils.GetRawUserInputP(s.conn, prompter, s.user.GetColorMode())
			default:
				panic("Unhandled case in switch statement (userInputMode)")
			}

			throttler.Sync()
			s.userInputChannel <- input
		}
	}()

	// Main loop
	for {
		input := s.getUserInputP(RawUserInput, s)
		if input == "" || input == "logout" || input == "quit" {
			return
		}

		if input == "." {
			s.WriteLine(s.lastInput)
			input = s.lastInput
		}

		s.lastInput = input

		if strings.HasPrefix(input, "/") {
			s.handleCommand(utils.Argify(input[1:]))
		} else {
			s.handleAction(utils.Argify(input))
		}
	}
}

func (s *Session) WriteLinef(line string, a ...interface{}) {
	s.WriteLineColor(color.White, line, a...)
}

func (s *Session) WriteLine(line string, a ...interface{}) {
	s.WriteLineColor(color.White, line, a...)
}

func (s *Session) WriteLineColor(color color.Color, line string, a ...interface{}) {
	s.printLine(color.Colorize(color, fmt.Sprintf(line, a...)))
}

func (s *Session) printLine(line string, a ...interface{}) {
	s.Write(fmt.Sprintf(line+"\r\n", a...))
}

func (s *Session) Write(text string) {
	s.user.Write(text)
}

func (s *Session) printError(err string, a ...interface{}) {
	s.WriteLineColor(color.Red, err, a...)
}

func (s *Session) clearLine() {
	utils.ClearLine(s.conn)
}

func (s *Session) asyncMessage(message string) {
	s.clearLine()
	s.WriteLine(message)
}

func (s *Session) GetInput(prompt string) string {
	return s.getUserInput(CleanUserInput, prompt)
}

func (s *Session) GetWindowSize() (int, int) {
	return s.user.GetWindowSize()
}

// getUserInput allows us to retrieve user input in a way that doesn't block the
// event loop by using channels and a separate Go routine to grab
// either the next user input or the next event.
func (s *Session) getUserInputP(inputMode userInputMode, prompter utils.Prompter) string {
	s.inputModeChannel <- inputMode
	s.prompterChannel <- prompter

	for {
		select {
		case input := <-s.userInputChannel:
			return input
		case event := <-s.eventChannel:
			if s.silentMode {
				continue
			}

			switch e := event.(type) {
			case events.TellEvent:
				s.replyId = e.From.GetId()
			case events.TickEvent:
				if !combat.InCombat(s.pc) {
					oldHps := s.pc.GetHitPoints()
					s.pc.Heal(5)
					newHps := s.pc.GetHitPoints()

					if oldHps != newHps {
						s.clearLine()
						s.Write(prompter.GetPrompt())
					}
				}
			}

			message := event.ToString(s.pc)
			if message != "" {
				s.asyncMessage(message)
				s.Write(prompter.GetPrompt())
			}

		case quitMessage := <-s.panicChannel:
			panic(quitMessage)
		}
	}
}

func (s *Session) getUserInput(inputMode userInputMode, prompt string) string {
	return s.getUserInputP(inputMode, utils.SimplePrompter(prompt))
}

func (s *Session) getCleanUserInput(prompt string) string {
	return s.getUserInput(CleanUserInput, prompt)
}

func (s *Session) getRawUserInput(prompt string) string {
	return s.getUserInput(RawUserInput, prompt)
}

func (s *Session) getConfirmation(prompt string) bool {
	answer := s.getCleanUserInput(color.Colorize(color.White, prompt))
	return answer == "y" || answer == "yes"
}

func (s *Session) getInt(prompt string, min, max int) (int, bool) {
	for {
		input := s.getRawUserInput(prompt)
		if input == "" {
			return 0, false
		}

		val, err := utils.Atoir(input, min, max)

		if err != nil {
			s.printError(err.Error())
		} else {
			return val, true
		}
	}
}

func (s *Session) GetPrompt() string {
	prompt := s.prompt
	prompt = strings.Replace(prompt, "%h", strconv.Itoa(s.pc.GetHitPoints()), -1)
	prompt = strings.Replace(prompt, "%H", strconv.Itoa(s.pc.GetHealth()), -1)

	if len(s.states) > 0 {
		states := make([]string, len(s.states))

		i := 0
		for key, value := range s.states {
			states[i] = fmt.Sprintf("%s:%s", key, value)
			i++
		}

		prompt = fmt.Sprintf("%s %s", states, prompt)
	}

	return color.Colorize(color.White, prompt)
}

func (s *Session) currentZone() types.Zone {
	return model.GetZone(s.GetRoom().GetZoneId())
}

func (s *Session) handleAction(action string, arg string) {
	if arg == "" {
		direction := types.StringToDirection(action)

		if direction != types.DirectionNone {
			if s.GetRoom().HasExit(direction) {
				err := model.MoveCharacter(s.pc, direction)
				if err == nil {
					s.PrintRoom()
				} else {
					s.printError(err.Error())
				}

			} else {
				s.printError("You can't go that way")
			}

			return
		}
	}

	handler, found := actions[action]

	if found {
		if handler.alias != "" {
			handler = actions[handler.alias]
		}
		handler.exec(s, arg)
	} else {
		s.printError("You can't do that")
	}
}

func (s *Session) handleCommand(name string, arg string) {
	if len(name) == 0 {
		return
	}

	if name[0] == '/' && s.user.IsAdmin() {
		quickRoom(s, name[1:])
		return
	}

	command, found := commands[name]

	if found {
		if command.alias != "" {
			command = commands[command.alias]
		}

		if command.admin && !s.user.IsAdmin() {
			s.printError("You don't have permission to do that")
		} else {
			command.exec(command, s, arg)
		}
	} else {
		s.printError("Unrecognized command: %s", name)
	}
}

func (s *Session) GetRoom() types.Room {
	return model.GetRoom(s.pc.GetRoomId())
}

func (s *Session) PrintRoom() {
	s.printRoom(s.GetRoom())
}

func (s *Session) printRoom(room types.Room) {
	pcs := model.PlayerCharactersIn(s.pc.GetRoomId(), s.pc.GetId())
	npcs := model.NpcsIn(room.GetId())
	items := model.ItemsIn(room.GetId())
	store := model.StoreIn(room.GetId())

	var area types.Area
	if room.GetAreaId() != nil {
		area = model.GetArea(room.GetAreaId())
	}

	var str string

	areaStr := ""
	if area != nil {
		areaStr = fmt.Sprintf("%s - ", area.GetName())
	}

	str = fmt.Sprintf("\r\n %v>>> %v%s%s %v<<< %v(%v %v %v)\r\n\r\n %v%s\r\n\r\n",
		color.White, color.Blue,
		areaStr, room.GetTitle(),
		color.White, color.Blue,
		room.GetLocation().X, room.GetLocation().Y, room.GetLocation().Z,
		color.White,
		room.GetDescription())

	if store != nil {
		str = fmt.Sprintf("%s Store: %s\r\n\r\n", str, color.Colorize(color.Blue, store.GetName()))
	}

	extraNewLine := ""

	if len(pcs) > 0 {
		str = fmt.Sprintf("%s %sAlso here:", str, color.Blue)

		names := make([]string, len(pcs))
		for i, char := range pcs {
			names[i] = color.Colorize(color.White, char.GetName())
		}
		str = fmt.Sprintf("%s %s \r\n", str, strings.Join(names, color.Colorize(color.Blue, ", ")))

		extraNewLine = "\r\n"
	}

	if len(npcs) > 0 {
		str = fmt.Sprintf("%s %s", str, color.Colorize(color.Blue, "NPCs: "))

		names := make([]string, len(npcs))
		for i, npc := range npcs {
			names[i] = color.Colorize(color.White, npc.GetName())
		}
		str = fmt.Sprintf("%s %s \r\n", str, strings.Join(names, color.Colorize(color.Blue, ", ")))

		extraNewLine = "\r\n"
	}

	if len(items) > 0 {
		itemMap := make(map[string]int)
		var nameList []string

		for _, item := range items {
			if item == nil {
				continue
			}

			_, found := itemMap[item.GetName()]
			if !found {
				nameList = append(nameList, item.GetName())
			}
			itemMap[item.GetName()]++
		}

		sort.Strings(nameList)

		str = str + " " + color.Colorize(color.Blue, "Items: ")

		var names []string
		for _, name := range nameList {
			if itemMap[name] > 1 {
				name = fmt.Sprintf("%s x%v", name, itemMap[name])
			}
			names = append(names, color.Colorize(color.White, name))
		}
		str = str + strings.Join(names, color.Colorize(color.Blue, ", ")) + "\r\n"

		extraNewLine = "\r\n"
	}

	str = str + extraNewLine + " " + color.Colorize(color.Blue, "Exits: ")

	var exitList []string
	for _, direction := range room.GetExits() {
		exitList = append(exitList, utils.DirectionToExitString(direction))
	}

	if len(exitList) == 0 {
		str = str + color.Colorize(color.White, "None")
	} else {
		str = str + strings.Join(exitList, " ")
	}

	if len(room.GetLinks()) > 0 {
		str = fmt.Sprintf("%s\r\n\r\n %s %s",
			str,
			color.Colorize(color.Blue, "Other exits:"),
			color.Colorize(color.White, strings.Join(room.LinkNames(), ", ")),
		)
	}

	str = str + "\r\n"

	s.WriteLine(str)
}

func (s *Session) execMenu(title string, build func(*utils.Menu)) {
	utils.ExecMenu(title, s, build)
}
