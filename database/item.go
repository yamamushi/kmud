package database

import (
	"github.com/yamamushi/kmud-2020/types"
	"github.com/yamamushi/kmud-2020/utils"
)

type Template struct {
	DbObject `bson:",inline"`
	Name     string
	Value    int
	Weight   int
	Capacity int
}

type Item struct {
	Container `bson:",inline"`

	TemplateId  types.Id
	Locked      bool
	ContainerId types.Id
}

func NewTemplate(name string) *Template {
	template := &Template{
		Name: utils.FormatName(name),
	}
	dbinit(template)
	return template
}

func NewItem(templateId types.Id) *Item {
	item := &Item{
		TemplateId: templateId,
	}
	dbinit(item)
	return item
}

// Template

func (t *Template) GetName() string {
	t.ReadLock()
	defer t.ReadUnlock()

	return t.Name
}

func (t *Template) SetName(name string) {
	t.writeLock(func() {
		t.Name = utils.FormatName(name)
	})
}

func (t *Template) SetValue(value int) {
	t.writeLock(func() {
		t.Value = value
	})
}

func (t *Template) GetValue() int {
	t.ReadLock()
	defer t.ReadUnlock()
	return t.Value
}

func (t *Template) GetWeight() int {
	t.ReadLock()
	defer t.ReadUnlock()
	return t.Weight
}

func (t *Template) SetWeight(weight int) {
	t.writeLock(func() {
		t.Weight = weight
	})
}

func (t *Template) GetCapacity() int {
	t.ReadLock()
	defer t.ReadUnlock()
	return t.Capacity
}

func (t *Template) SetCapacity(capacity int) {
	t.writeLock(func() {
		t.Capacity = capacity
	})
}

// Item

func (i *Item) GetTemplateId() types.Id {
	i.ReadLock()
	defer i.ReadUnlock()
	return i.TemplateId
}

func (i *Item) GetTemplate() types.Template {
	i.ReadLock()
	defer i.ReadUnlock()
	return Retrieve(i.TemplateId, types.TemplateType).(types.Template)
}

func (i *Item) GetName() string {
	return i.GetTemplate().GetName()
}

func (i *Item) GetValue() int {
	return i.GetTemplate().GetValue()
}

func (i *Item) GetCapacity() int {
	return i.GetTemplate().GetCapacity()
}

func (i *Item) IsLocked() bool {
	i.ReadLock()
	defer i.ReadUnlock()

	return i.Locked
}

func (i *Item) SetLocked(locked bool) {
	i.writeLock(func() {
		i.Locked = locked
	})
}

func (i *Item) GetContainerId() types.Id {
	i.ReadLock()
	defer i.ReadUnlock()
	return i.ContainerId
}

func (i *Item) SetContainerId(id types.Id, from types.Id) bool {
	i.WriteLock()
	if from != i.ContainerId {
		i.WriteUnlock()
		return false
	}
	i.ContainerId = id
	i.WriteUnlock()
	i.syncModified()
	return true
}
