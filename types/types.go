package types

import (
	"github.com/yamamushi/kmud-2020/color"
	"net"

	"github.com/yamamushi/kmud-2020/utils/naturalsort"
)

type Id interface {
	String() string
	Hex() string
}

type ObjectType string

const (
	NpcType      ObjectType = "Npc"
	PcType       ObjectType = "Pc"
	SpawnerType  ObjectType = "Spawner"
	UserType     ObjectType = "User"
	ZoneType     ObjectType = "Zone"
	AreaType     ObjectType = "Area"
	RoomType     ObjectType = "Room"
	TemplateType ObjectType = "Template"
	ItemType     ObjectType = "Item"
	SkillType    ObjectType = "Skill"
	EffectType   ObjectType = "Effect"
	StoreType    ObjectType = "Store"
	WorldType    ObjectType = "World"
)

type Identifiable interface {
	GetId() Id
}

type ReadLockable interface {
	ReadLock()
	ReadUnlock()
}

type Destroyable interface {
	Destroy()
	IsDestroyed() bool
}

type Locateable interface {
	GetRoomId() Id
}

type Nameable interface {
	GetName() string
	SetName(string)
}

type Loginable interface {
	IsOnline() bool
	SetOnline(bool)
}

type Container interface {
	AddCash(int)
	RemoveCash(int)
	GetCash() int
	SetCapacity(int)
	GetCapacity() int
}

type Object interface {
	Identifiable
	ReadLockable
	Destroyable
	SetId(Id)
}

type PC interface {
	Character
	Loginable
}

type PCList []PC

func (l PCList) Characters() CharacterList {
	chars := make(CharacterList, len(l))
	for i, pc := range l {
		chars[i] = pc
	}
	return chars
}

type NPC interface {
	Character
	SetRoaming(bool)
	GetRoaming() bool
	SetConversation(string)
	GetConversation() string
	PrettyConversation() string
}

type NPCList []NPC

type Spawner interface {
	Character
	GetAreaId() Id
	SetCount(int)
	GetCount() int
}

type SpawnerList []Spawner

func (l NPCList) Characters() CharacterList {
	chars := make(CharacterList, len(l))
	for i, npc := range l {
		chars[i] = npc
	}
	return chars
}

type Room interface {
	Object
	Container
	GetZoneId() Id
	GetAreaId() Id
	SetAreaId(Id)
	GetLocation() Coordinate
	SetExitEnabled(Direction, bool)
	HasExit(Direction) bool
	NextLocation(Direction) Coordinate
	GetExits() []Direction
	GetTitle() string
	SetTitle(string)
	GetDescription() string
	SetDescription(string)
	SetLink(string, Id)
	RemoveLink(string)
	GetLinks() map[string]Id
	LinkNames() []string
	SetLocked(Direction, bool)
	IsLocked(Direction) bool
}

type RoomList []Room

type Area interface {
	Object
	Nameable
}

type AreaList []Area

type Zone interface {
	Object
	Nameable
}

type ZoneList []Zone

type Time interface {
	String() string
}

type World interface {
	GetTime() Time
	AdvanceTime()
}

type Communicable interface {
	WriteLine(string, ...interface{})
	Write(string)
	GetInput(prompt string) string
	GetWindowSize() (int, int)
}

type User interface {
	Object
	Nameable
	Loginable
	Communicable
	VerifyPassword(string) bool
	SetConnection(net.Conn)
	GetConnection() net.Conn
	SetWindowSize(int, int)
	SetTerminalType(string)
	GetTerminalType() string
	GetColorMode() color.ColorMode
	SetColorMode(color.ColorMode)
	IsAdmin() bool
	SetAdmin(bool)
}

type UserList []User

func (l UserList) Len() int {
	return len(l)
}

func (l UserList) Less(i, j int) bool {
	return naturalsort.NaturalLessThan(l[i].GetName(), l[j].GetName())
}

func (l UserList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Template interface {
	Object
	Nameable
	SetValue(int)
	GetValue() int
	SetWeight(int)
	GetWeight() int
	GetCapacity() int
	SetCapacity(int)
}

type TemplateList []Template

func (l TemplateList) Names() []string {
	names := make([]string, len(l))
	for i, item := range l {
		names[i] = item.GetName()
	}
	return names
}

func (l TemplateList) Len() int {
	return len(l)
}

func (l TemplateList) Less(i, j int) bool {
	return naturalsort.NaturalLessThan(l[i].GetName(), l[j].GetName())
}

func (l TemplateList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Item interface {
	Object
	Container
	GetTemplateId() Id
	GetName() string
	GetValue() int
	SetLocked(bool)
	IsLocked() bool
	GetContainerId() Id
	SetContainerId(Id, Id) bool
}

type ItemList []Item

func (l ItemList) Names() []string {
	names := make([]string, len(l))
	for i, item := range l {
		names[i] = item.GetName()
	}
	return names
}

func (l ItemList) Len() int {
	return len(l)
}

func (l ItemList) Less(i, j int) bool {
	return naturalsort.NaturalLessThan(l[i].GetName(), l[j].GetName())
}

func (l ItemList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type Skill interface {
	Object
	Nameable
	GetEffects() []Id
	AddEffect(Id)
	RemoveEffect(Id)
	HasEffect(Id) bool
}

type SkillList []Skill

type EffectKind string

const (
	HitpointEffect EffectKind = "hitpoint"
	SilenceEffect  EffectKind = "silence"
	StunEffect     EffectKind = "stun"
)

func (l SkillList) Names() []string {
	names := make([]string, len(l))
	for i, skill := range l {
		names[i] = skill.GetName()
	}
	return names
}

type Effect interface {
	Object
	Nameable
	SetPower(int)
	GetPower() int
	SetCost(int)
	GetCost() int
	GetType() EffectKind
	SetType(EffectKind)
	GetVariance() int
	SetVariance(int)
	GetSpeed() int
	SetSpeed(int)
	GetTime() int
	SetTime(int)
}

type EffectList []Effect

func (l EffectList) Names() []string {
	names := make([]string, len(l))
	for i, skill := range l {
		names[i] = skill.GetName()
	}
	return names
}

type Store interface {
	Object
	Nameable
	Container
}

type Purchaser interface {
	GetId() Id
	AddCash(int)
	RemoveCash(int)
}
