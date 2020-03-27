package olddatabase

import (
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type Skill struct {
	DbObject `bson:",inline"`

	Effects utils.Set
	Name    string
}

func NewSkill(name string) *Skill {
	skill := &Skill{
		Name: utils.FormatName(name),
	}

	dbinit(skill)
	return skill
}

func (s *Skill) GetName() string {
	s.ReadLock()
	defer s.ReadUnlock()
	return s.Name
}

func (s *Skill) SetName(name string) {
	s.writeLock(func() {
		s.Name = utils.FormatName(name)
	})
}

func (s *Skill) AddEffect(id types.Id) {
	s.writeLock(func() {
		if s.Effects == nil {
			s.Effects = utils.Set{}
		}
		s.Effects.Insert(id.Hex())
	})
}

func (s *Skill) RemoveEffect(id types.Id) {
	s.writeLock(func() {
		s.Effects.Remove(id.Hex())
	})
}

func (s *Skill) GetEffects() []types.Id {
	s.ReadLock()
	defer s.ReadUnlock()
	return idSetToList(s.Effects)
}

func (s *Skill) HasEffect(id types.Id) bool {
	s.ReadLock()
	defer s.ReadUnlock()
	return s.Effects.Contains(id.Hex())
}
