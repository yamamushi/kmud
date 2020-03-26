package database

import (
	"github.com/yamamushi/kmud/types"
	"github.com/yamamushi/kmud/utils"
)

type Store struct {
	Container `bson:",inline"`

	Name   string
	RoomId types.Id
}

func NewStore(name string, roomId types.Id) *Store {
	store := &Store{
		Name:   utils.FormatName(name),
		RoomId: roomId,
	}

	dbinit(store)
	return store
}

func (s *Store) GetName() string {
	s.ReadLock()
	defer s.ReadUnlock()

	return s.Name
}

func (s *Store) SetName(name string) {
	s.writeLock(func() {
		s.Name = utils.FormatName(name)
	})
}
