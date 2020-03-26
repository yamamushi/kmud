package database

import (
	"sync"

	"github.com/yamamushi/kmud/types"
	"github.com/yamamushi/kmud/utils"
	"gopkg.in/mgo.v2/bson"
)

type DbObject struct {
	Id types.Id `bson:"_id"`

	mutex     sync.RWMutex
	destroyed bool
}

func (d *DbObject) SetId(id types.Id) {
	d.Id = id
}

func (d *DbObject) GetId() types.Id {
	return d.Id
}

func (d *DbObject) ReadLock() {
	d.mutex.RLock()
}

func (d *DbObject) ReadUnlock() {
	d.mutex.RUnlock()
}

func (d *DbObject) WriteLock() {
	d.mutex.Lock()
}

func (d *DbObject) writeLock(worker func()) {
	d.WriteLock()
	defer d.WriteUnlock()
	defer d.modified()
	worker()
}

func (d *DbObject) WriteUnlock() {
	d.mutex.Unlock()
}

func (d *DbObject) Destroy() {
	d.WriteLock()
	defer d.WriteUnlock()

	d.destroyed = true
}

func (d *DbObject) IsDestroyed() bool {
	d.ReadLock()
	defer d.ReadUnlock()

	return d.destroyed
}

func (d *DbObject) modified() {
	modifiedObjectChannel <- d.Id
}

func (d *DbObject) syncModified() {
	commitObject(d.Id)
}

func idSetToList(set utils.Set) []types.Id {
	ids := make([]types.Id, len(set))

	i := 0
	for id := range set {
		ids[i] = bson.ObjectIdHex(id)
		i++
	}

	return ids
}
