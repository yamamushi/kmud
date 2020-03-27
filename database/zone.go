package database

import "github.com/yamamushi/kmud-2020/utils"

type Zone struct {
	DbObject `bson:",inline"`

	Name string
}

func NewZone(name string) *Zone {
	zone := &Zone{
		Name: utils.FormatName(name),
	}

	dbinit(zone)
	return zone
}

func (z *Zone) GetName() string {
	z.ReadLock()
	defer z.ReadUnlock()
	return z.Name
}

func (z *Zone) SetName(name string) {
	z.writeLock(func() {
		z.Name = utils.FormatName(name)
	})
}
