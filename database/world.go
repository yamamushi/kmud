package database

import (
	"fmt"
	"time"

	"github.com/yamamushi/kmud/types"
)

type World struct {
	DbObject `bson:",inline"`
	Time     time.Time
}

func NewWorld() *World {
	world := &World{Time: time.Now()}

	dbinit(world)
	return world
}

type _time struct {
	hour int
	min  int
	sec  int
}

func (t _time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.hour, t.min, t.sec)
}

const _TIME_MULTIPLIER = 3

func (w *World) GetTime() types.Time {
	w.ReadLock()
	defer w.ReadUnlock()

	hour, min, sec := w.Time.Clock()
	return _time{hour: hour, min: min, sec: sec}
}

func (w *World) AdvanceTime() {
	w.writeLock(func() {
		w.Time = w.Time.Add(3 * time.Second)
	})
}
