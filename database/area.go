package database

import (
	"github.com/yamamushi/kmud/types"
	"github.com/yamamushi/kmud/utils"
)

type Area struct {
	DbObject `bson:",inline"`
	Name     string
	ZoneId   types.Id
}

func NewArea(name string, zone types.Id) *Area {
	area := &Area{
		ZoneId: zone,
		Name:   utils.FormatName(name),
	}
	dbinit(area)
	return area
}

func (a *Area) GetName() string {
	a.ReadLock()
	defer a.ReadUnlock()
	return a.Name
}

func (a *Area) SetName(name string) {
	a.writeLock(func() {
		a.Name = name
	})
}

func (a *Area) GetZoneId() types.Id {
	a.ReadLock()
	defer a.ReadUnlock()
	return a.ZoneId
}
