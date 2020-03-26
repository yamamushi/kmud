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

func (self *Area) GetName() string {
	self.ReadLock()
	defer self.ReadUnlock()
	return self.Name
}

func (self *Area) SetName(name string) {
	self.writeLock(func() {
		self.Name = name
	})
}

func (self *Area) GetZoneId() types.Id {
	self.ReadLock()
	defer self.ReadUnlock()
	return self.ZoneId
}
