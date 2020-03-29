package olddatabase

import (
	"fmt"
	"github.com/yamamushi/kmud-2020/color"

	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type Character struct {
	Container `bson:",inline"`

	RoomId    types.Id `bson:",omitempty"`
	Name      string
	HitPoints int
	Skills    utils.Set

	Strength int
	Vitality int
}

type Pc struct {
	Character `bson:",inline"`

	UserId types.Id
	online bool
}

type Npc struct {
	Character `bson:",inline"`

	SpawnerId types.Id `bson:",omitempty"`

	Roaming      bool
	Conversation string
}

type Spawner struct {
	Character `bson:",inline"`

	AreaId types.Id
	Count  int
}

func NewPc(name string, userId types.Id, roomId types.Id) *Pc {
	pc := &Pc{
		UserId: userId,
		online: false,
	}

	pc.initCharacter(name, types.PcType, roomId)
	dbinit(pc)
	return pc
}

func NewNpc(name string, roomId types.Id, spawnerId types.Id) *Npc {
	npc := &Npc{
		SpawnerId: spawnerId,
	}

	npc.initCharacter(name, types.NpcType, roomId)
	dbinit(npc)
	return npc
}

func NewSpawner(name string, areaId types.Id) *Spawner {
	spawner := &Spawner{
		AreaId: areaId,
		Count:  1,
	}

	spawner.initCharacter(name, types.SpawnerType, nil)
	dbinit(spawner)
	return spawner
}

func (c *Character) initCharacter(name string, objType types.ObjectType, roomId types.Id) {
	c.RoomId = roomId
	c.Cash = 0
	c.HitPoints = 100
	c.Name = utils.FormatName(name)

	c.Strength = 10
	c.Vitality = 100
}

func (c *Character) GetName() string {
	c.ReadLock()
	defer c.ReadUnlock()

	return c.Name
}

func (c *Character) SetName(name string) {
	c.writeLock(func() {
		c.Name = utils.FormatName(name)
	})
}

func (c *Character) GetCapacity() int {
	return c.GetStrength() * 10
}

func (c *Character) GetStrength() int {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.Strength
}

func (pc *Pc) SetOnline(online bool) {
	pc.WriteLock()
	defer pc.WriteUnlock()
	pc.online = online
}

func (pc *Pc) IsOnline() bool {
	pc.ReadLock()
	defer pc.ReadUnlock()
	return pc.online
}

func (c *Character) SetRoomId(id types.Id) {
	c.writeLock(func() {
		c.RoomId = id
	})
}

func (c *Character) GetRoomId() types.Id {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.RoomId
}

func (pc *Pc) SetUserId(id types.Id) {
	pc.writeLock(func() {
		pc.UserId = id
	})
}

func (pc *Pc) GetUserId() types.Id {
	pc.ReadLock()
	defer pc.ReadUnlock()
	return pc.UserId
}

func (c *Character) AddSkill(id types.Id) {
	c.writeLock(func() {
		if c.Skills == nil {
			c.Skills = utils.Set{}
		}
		c.Skills.Insert(id.Hex())
	})
}

func (c *Character) RemoveSkill(id types.Id) {
	c.writeLock(func() {
		c.Skills.Remove(id.Hex())
	})
}

func (c *Character) HasSkill(id types.Id) bool {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.Skills.Contains(id.Hex())
}

func (c *Character) GetSkills() []types.Id {
	c.ReadLock()
	defer c.ReadUnlock()
	return idSetToList(c.Skills)
}

func (n *Npc) SetConversation(conversation string) {
	n.writeLock(func() {
		n.Conversation = conversation
	})
}

func (n *Npc) GetConversation() string {
	n.ReadLock()
	defer n.ReadUnlock()
	return n.Conversation
}

func (n *Npc) PrettyConversation() string {
	conv := n.GetConversation()

	if conv == "" {
		return fmt.Sprintf("%n has nothing to say", n.GetName())
	}

	return fmt.Sprintf("%n%n",
		color.Colorize(color.Blue, n.GetName()),
		color.Colorize(color.White, ": "+conv))
}

func (c *Character) SetHealth(health int) {
	c.writeLock(func() {
		c.Vitality = health
		if c.HitPoints > c.Vitality {
			c.HitPoints = c.Vitality
		}
	})
}

func (c *Character) GetHealth() int {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.Vitality
}

func (c *Character) SetHitPoints(hitpoints int) {
	c.writeLock(func() {
		if hitpoints > c.Vitality {
			hitpoints = c.Vitality
		}
		c.HitPoints = hitpoints
	})
}

func (c *Character) GetHitPoints() int {
	c.ReadLock()
	defer c.ReadUnlock()
	return c.HitPoints
}

func (c *Character) Hit(hitpoints int) {
	c.SetHitPoints(c.GetHitPoints() - hitpoints)
}

func (c *Character) Heal(hitpoints int) {
	c.SetHitPoints(c.GetHitPoints() + hitpoints)
}

func (n *Npc) GetRoaming() bool {
	n.ReadLock()
	defer n.ReadUnlock()
	return n.Roaming
}

func (n *Npc) SetRoaming(roaming bool) {
	n.writeLock(func() {
		n.Roaming = roaming
	})
}

func (s *Spawner) SetCount(count int) {
	s.writeLock(func() {
		s.Count = count
	})
}

func (s *Spawner) GetCount() int {
	s.ReadLock()
	defer s.ReadUnlock()
	return s.Count
}

func (s *Spawner) GetAreaId() types.Id {
	s.ReadLock()
	defer s.ReadUnlock()
	return s.AreaId
}
