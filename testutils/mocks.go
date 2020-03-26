package testutils

import (
	"github.com/yamamushi/kmud/types"
	"gopkg.in/mgo.v2/bson"
)

type MockId string

func (m MockId) String() string {
	return string(m)
}

func (m MockId) Hex() string {
	return string(m)
}

type MockIdentifiable struct {
	Id types.Id
}

func (m MockIdentifiable) GetId() types.Id {
	return m.Id
}

func (m MockObject) SetId(types.Id) {
}

type MockNameable struct {
	Name string
}

func (n MockNameable) GetName() string {
	return n.Name
}

func (s *MockNameable) SetName(name string) {
	s.Name = name
}

type MockDestroyable struct {
}

func (MockDestroyable) Destroy() {
}

func (d MockDestroyable) IsDestroyed() bool {
	return false
}

type MockReadLocker struct {
}

func (*MockReadLocker) ReadLock() {
}

func (*MockReadLocker) ReadUnlock() {
}

type MockObject struct {
	MockIdentifiable
	MockReadLocker
	MockDestroyable
}

type MockZone struct {
	MockIdentifiable
}

func NewMockZone() *MockZone {
	return &MockZone{
		MockIdentifiable{Id: bson.NewObjectId()},
	}
}

type MockRoom struct {
	MockIdentifiable
}

func NewMockRoom() *MockRoom {
	return &MockRoom{
		MockIdentifiable{Id: bson.NewObjectId()},
	}
}

type MockUser struct {
	MockIdentifiable
}

func NewMockUser() *MockUser {
	return &MockUser{
		MockIdentifiable{Id: bson.NewObjectId()},
	}
}

type MockContainer struct {
}

func (*MockContainer) AddCash(int) {
}

func (*MockContainer) GetCash() int {
	return 0
}

func (*MockContainer) RemoveCash(int) {
}

func (*MockContainer) AddItem(types.Id) {
}

func (*MockContainer) RemoveItem(types.Id) bool {
	return true
}

func (*MockContainer) GetCapacity() int {
	return 0
}

func (*MockContainer) SetCapacity(int) {
}

func (*MockContainer) GetItems() types.ItemList {
	return types.ItemList{}
}

type MockCharacter struct {
	MockObject
	MockNameable
	MockContainer
}

func (*MockCharacter) GetHealth() int {
	return 1
}

func (*MockCharacter) SetHealth(int) {
}

func (*MockCharacter) GetHitPoints() int {
	return 1
}

func (*MockCharacter) Heal(int) {
}

func (*MockCharacter) Hit(int) {
}

func (*MockCharacter) SetHitPoints(int) {
}

func (*MockCharacter) GetWeight() int {
	return 0
}

type MockPC struct {
	MockCharacter
	RoomId types.Id
}

func NewMockPC() *MockPC {
	return &MockPC{
		MockCharacter: MockCharacter{
			MockObject: MockObject{
				MockIdentifiable: MockIdentifiable{Id: bson.NewObjectId()},
			},
			MockNameable: MockNameable{Name: "Mock PC"},
		},
		RoomId: bson.NewObjectId(),
	}
}

func (p MockPC) GetRoomId() types.Id {
	return p.RoomId
}

func (p MockPC) IsOnline() bool {
	return true
}

func (p MockPC) SetRoomId(types.Id) {
}

func (p MockPC) GetSkills() []types.Id {
	return []types.Id{}
}

func (p MockPC) AddSkill(types.Id) {
}
