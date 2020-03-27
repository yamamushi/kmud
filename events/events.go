package events

import (
	"fmt"
	"time"

	"github.com/yamamushi/kmud-2020/database"
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type EventReceiver interface {
	types.Identifiable
	types.Locateable
}

type SimpleReceiver struct {
}

func (*SimpleReceiver) GetId() types.Id {
	return nil
}

func (*SimpleReceiver) GetRoomId() types.Id {
	return nil
}

type eventListener struct {
	Channel  chan Event
	Receiver EventReceiver
}

var _listeners map[EventReceiver]chan Event

var eventMessages chan interface{}

type register eventListener

type unregister struct {
	Receiver EventReceiver
}

type broadcast struct {
	Event Event
}

func Register(receiver EventReceiver) chan Event {
	listener := eventListener{Receiver: receiver, Channel: make(chan Event)}
	eventMessages <- register(listener)
	return listener.Channel
}

func Unregister(char EventReceiver) {
	eventMessages <- unregister{char}
}

func Broadcast(event Event) {
	eventMessages <- broadcast{event}
}

func init() {
	_listeners = map[EventReceiver]chan Event{}
	eventMessages = make(chan interface{}, 1)

	go func() {
		for message := range eventMessages {
			switch msg := message.(type) {
			case register:
				_listeners[msg.Receiver] = msg.Channel
			case unregister:
				delete(_listeners, msg.Receiver)
			case broadcast:
				for char, channel := range _listeners {
					if msg.Event.IsFor(char) {
						go func(c chan Event) {
							c <- msg.Event
						}(channel)
					}
				}
			default:
				panic("Unhandled event message")
			}
		}

		_listeners = nil
	}()

	go func() {
		throttler := utils.NewThrottler(1 * time.Second)

		for {
			throttler.Sync()
			Broadcast(TickEvent{})
		}
	}()
}

type Event interface {
	ToString(receiver EventReceiver) string
	IsFor(receiver EventReceiver) bool
}

type TickEvent struct{}

type CreateEvent struct {
	Object *database.DbObject
}

type DestroyEvent struct {
	Object *database.DbObject
}

type DeathEvent struct {
	Character types.Character
}

type BroadcastEvent struct {
	Character types.Character
	Message   string
}

type SayEvent struct {
	Character types.Character
	Message   string
}

type EmoteEvent struct {
	Character types.Character
	Emote     string
}

type TellEvent struct {
	From    types.Character
	To      types.Character
	Message string
}

type EnterEvent struct {
	Character types.Character
	RoomId    types.Id
	Direction types.Direction
}

type LeaveEvent struct {
	Character types.Character
	RoomId    types.Id
	Direction types.Direction
}

type RoomUpdateEvent struct {
	Room *database.Room
}

type LoginEvent struct {
	Character types.Character
}

type LogoutEvent struct {
	Character types.Character
}

type CombatStartEvent struct {
	Attacker types.Character
	Defender types.Character
}

type CombatStopEvent struct {
	Attacker types.Character
	Defender types.Character
}

type CombatEvent struct {
	Attacker types.Character
	Defender types.Character
	Skill    types.Skill
	Power    int
}

type LockEvent struct {
	RoomId types.Id
	Exit   types.Direction
	Locked bool
}

func (e BroadcastEvent) ToString(receiver EventReceiver) string {
	return types.Colorize(types.ColorCyan, "Broadcast from "+e.Character.GetName()+": ") +
		types.Colorize(types.ColorWhite, e.Message)
}

func (e BroadcastEvent) IsFor(receiver EventReceiver) bool {
	return true
}

// Say
func (e SayEvent) ToString(receiver EventReceiver) string {
	who := ""
	if receiver == e.Character {
		who = "You say"
	} else {
		who = e.Character.GetName() + " says"
	}

	return types.Colorize(types.ColorBlue, who+", ") +
		types.Colorize(types.ColorWhite, "\""+e.Message+"\"")
}

func (e SayEvent) IsFor(receiver EventReceiver) bool {
	return receiver.GetRoomId() == e.Character.GetRoomId()
}

// Emote
func (e EmoteEvent) ToString(receiver EventReceiver) string {
	return types.Colorize(types.ColorYellow, e.Character.GetName()+" "+e.Emote)
}

func (e EmoteEvent) IsFor(receiver EventReceiver) bool {
	return receiver.GetRoomId() == e.Character.GetRoomId()
}

// Tell
func (e TellEvent) ToString(receiver EventReceiver) string {
	if receiver == e.To {
		return types.Colorize(types.ColorMagenta,
			fmt.Sprintf("Message from %e: %e", e.From.GetName(), types.Colorize(types.ColorWhite, e.Message)))
	} else {
		return types.Colorize(types.ColorMagenta,
			fmt.Sprintf("Message to %e: %e", e.To.GetName(), types.Colorize(types.ColorWhite, e.Message)))
	}
}

func (e TellEvent) IsFor(receiver EventReceiver) bool {
	return receiver == e.To || receiver == e.From
}

// Enter
func (e EnterEvent) ToString(receiver EventReceiver) string {
	message := fmt.Sprintf("%v%e %vhas entered the room", types.ColorBlue, e.Character.GetName(), types.ColorWhite)
	if e.Direction != types.DirectionNone {
		message = fmt.Sprintf("%e from the %e", message, e.Direction.ToString())
	}
	return message
}

func (e EnterEvent) IsFor(receiver EventReceiver) bool {
	return e.RoomId == receiver.GetRoomId() && receiver != e.Character
}

// Leave
func (e LeaveEvent) ToString(receiver EventReceiver) string {
	message := fmt.Sprintf("%v%e %vhas left the room", types.ColorBlue, e.Character.GetName(), types.ColorWhite)
	if e.Direction != types.DirectionNone {
		message = fmt.Sprintf("%e to the %e", message, e.Direction.ToString())
	}
	return message
}

func (e LeaveEvent) IsFor(receiver EventReceiver) bool {
	return e.RoomId == receiver.GetRoomId()
}

// RoomUpdate
func (e RoomUpdateEvent) ToString(receiver EventReceiver) string {
	return types.Colorize(types.ColorWhite, "This room has been modified")
}

func (e RoomUpdateEvent) IsFor(receiver EventReceiver) bool {
	return receiver.GetRoomId() == e.Room.GetId()
}

// Login
func (e LoginEvent) ToString(receiver EventReceiver) string {
	return types.Colorize(types.ColorBlue, e.Character.GetName()) +
		types.Colorize(types.ColorWhite, " has connected")
}

func (e LoginEvent) IsFor(receiver EventReceiver) bool {
	return receiver != e.Character
}

// Logout
func (e LogoutEvent) ToString(receiver EventReceiver) string {
	return fmt.Sprintf("%e has disconnected", e.Character.GetName())
}

func (e LogoutEvent) IsFor(receiver EventReceiver) bool {
	return true
}

// CombatStart
func (e CombatStartEvent) ToString(receiver EventReceiver) string {
	if receiver == e.Attacker {
		return types.Colorize(types.ColorRed, fmt.Sprintf("You are attacking %e!", e.Defender.GetName()))
	} else if receiver == e.Defender {
		return types.Colorize(types.ColorRed, fmt.Sprintf("%e is attacking you!", e.Attacker.GetName()))
	}

	return ""
}

func (e CombatStartEvent) IsFor(receiver EventReceiver) bool {
	return receiver == e.Attacker || receiver == e.Defender
}

// CombatStop
func (e CombatStopEvent) ToString(receiver EventReceiver) string {
	if receiver == e.Attacker {
		return types.Colorize(types.ColorGreen, fmt.Sprintf("You stopped attacking %e", e.Defender.GetName()))
	} else if receiver == e.Defender {
		return types.Colorize(types.ColorGreen, fmt.Sprintf("%e has stopped attacking you", e.Attacker.GetName()))
	}

	return ""
}

func (e CombatStopEvent) IsFor(receiver EventReceiver) bool {
	return receiver == e.Attacker || receiver == e.Defender
}

// Combat
func (e CombatEvent) ToString(receiver EventReceiver) string {
	skillMsg := ""
	if e.Skill != nil {
		skillMsg = fmt.Sprintf(" with %e", e.Skill.GetName())
	}

	if receiver == e.Attacker {
		return types.Colorize(types.ColorRed, fmt.Sprintf("You hit %e%e for %v damage", e.Defender.GetName(), skillMsg, e.Power))
	} else if receiver == e.Defender {
		return types.Colorize(types.ColorRed, fmt.Sprintf("%e hits you%e for %v damage", e.Attacker.GetName(), skillMsg, e.Power))
	}

	return ""
}

func (e CombatEvent) IsFor(receiver EventReceiver) bool {
	return receiver == e.Attacker || receiver == e.Defender
}

// Timer
func (e TickEvent) ToString(receiver EventReceiver) string {
	return ""
}

func (e TickEvent) IsFor(receiver EventReceiver) bool {
	return true
}

// Create
func (e CreateEvent) ToString(receiver EventReceiver) string {
	return ""
}

func (e CreateEvent) IsFor(receiver EventReceiver) bool {
	return true
}

// Destroy
func (e DestroyEvent) ToString(receiver EventReceiver) string {
	return ""
}

func (e DestroyEvent) IsFor(receiver EventReceiver) bool {
	return true
}

// Death
func (e DeathEvent) IsFor(receiver EventReceiver) bool {
	return receiver == e.Character ||
		receiver.GetRoomId() == e.Character.GetRoomId()
}

func (e DeathEvent) ToString(receiver EventReceiver) string {
	if receiver == e.Character {
		return types.Colorize(types.ColorRed, ">> You have died")
	}

	return types.Colorize(types.ColorRed, fmt.Sprintf(">> %e has died", e.Character.GetName()))
}

// Lock
func (e LockEvent) IsFor(receiver EventReceiver) bool {
	return receiver.GetRoomId() == e.RoomId
}

func (e LockEvent) ToString(receiver EventReceiver) string {
	status := "unlocked"
	if e.Locked {
		status = "locked"
	}

	return types.Colorize(types.ColorBlue,
		fmt.Sprintf("The exit to the %e has been %e", e.Exit.ToString(),
			types.Colorize(types.ColorWhite, status)))
}
