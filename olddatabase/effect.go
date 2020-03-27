package olddatabase

import (
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type Effect struct {
	DbObject `bson:",inline"`

	Type     types.EffectKind
	Name     string
	Power    int
	Cost     int
	Variance int
	Speed    int
	Time     int
}

func NewEffect(name string) types.Effect {
	effect := &Effect{
		Name:     utils.FormatName(name),
		Power:    1,
		Type:     types.HitpointEffect,
		Variance: 0,
		Time:     1,
	}

	dbinit(effect)
	return effect
}

func (e *Effect) GetName() string {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Name
}

func (e *Effect) SetName(name string) {
	e.writeLock(func() {
		e.Name = utils.FormatName(name)
	})
}

func (e *Effect) GetType() types.EffectKind {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Type
}

func (e *Effect) SetType(effectKind types.EffectKind) {
	e.writeLock(func() {
		e.Type = effectKind
	})
}

func (e *Effect) GetPower() int {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Power
}

func (e *Effect) SetPower(power int) {
	e.writeLock(func() {
		e.Power = power
	})
}

func (e *Effect) GetCost() int {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Cost
}

func (e *Effect) SetCost(cost int) {
	e.writeLock(func() {
		e.Cost = cost
	})
}

func (e *Effect) GetVariance() int {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Variance
}

func (e *Effect) SetVariance(variance int) {
	e.writeLock(func() {
		e.Variance = variance
	})
}

func (e *Effect) GetSpeed() int {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Speed
}

func (e *Effect) SetSpeed(speed int) {
	e.writeLock(func() {
		e.Speed = speed
	})
}

func (e *Effect) GetTime() int {
	e.ReadLock()
	defer e.ReadUnlock()
	return e.Time
}

func (e *Effect) SetTime(speed int) {
	e.writeLock(func() {
		e.Time = speed
	})
}
