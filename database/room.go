package database

import "github.com/yamamushi/kmud-2020/types"

type Exit struct {
	Locked bool
}

type Room struct {
	Container `bson:",inline"`

	ZoneId      types.Id
	AreaId      types.Id `bson:",omitempty"`
	Title       string
	Description string
	Links       map[string]types.Id
	Location    types.Coordinate

	Exits map[types.Direction]*Exit
}

func NewRoom(zoneId types.Id, location types.Coordinate) *Room {
	room := &Room{
		Title: "The Void",
		Description: "You are floating in the blackness of space. Complete darkness surrounds " +
			"you in all directions. There is no escape, there is no hope, just the emptiness. " +
			"You are likely to be eaten by a grue.",
		Location: location,
		ZoneId:   zoneId,
	}

	dbinit(room)
	return room
}

func (r *Room) HasExit(dir types.Direction) bool {
	r.ReadLock()
	defer r.ReadUnlock()

	_, found := r.Exits[dir]
	return found
}

func (r *Room) SetExitEnabled(dir types.Direction, enabled bool) {
	r.writeLock(func() {
		if r.Exits == nil {
			r.Exits = map[types.Direction]*Exit{}
		}
		if enabled {
			r.Exits[dir] = &Exit{}
		} else {
			delete(r.Exits, dir)
		}
	})
}

func (r *Room) SetLink(name string, roomId types.Id) {
	r.writeLock(func() {
		if r.Links == nil {
			r.Links = map[string]types.Id{}
		}
		r.Links[name] = roomId
	})
}

func (r *Room) RemoveLink(name string) {
	r.writeLock(func() {
		delete(r.Links, name)
	})
}

func (r *Room) GetLinks() map[string]types.Id {
	r.ReadLock()
	defer r.ReadUnlock()
	return r.Links
}

func (r *Room) LinkNames() []string {
	names := make([]string, len(r.GetLinks()))

	i := 0
	for name := range r.Links {
		names[i] = name
		i++
	}
	return names
}

func (r *Room) SetTitle(title string) {
	r.writeLock(func() {
		r.Title = title
	})
}

func (r *Room) GetTitle() string {
	r.ReadLock()
	defer r.ReadUnlock()
	return r.Title
}

func (r *Room) SetDescription(description string) {
	r.writeLock(func() {
		r.Description = description
	})
}

func (r *Room) GetDescription() string {
	r.ReadLock()
	defer r.ReadUnlock()
	return r.Description
}

func (r *Room) SetLocation(location types.Coordinate) {
	r.writeLock(func() {
		r.Location = location
	})
}

func (r *Room) GetLocation() types.Coordinate {
	r.ReadLock()
	defer r.ReadUnlock()
	return r.Location
}

func (r *Room) SetZoneId(zoneId types.Id) {
	r.writeLock(func() {
		r.ZoneId = zoneId
	})
}

func (r *Room) GetZoneId() types.Id {
	r.ReadLock()
	defer r.ReadUnlock()
	return r.ZoneId
}

func (r *Room) SetAreaId(areaId types.Id) {
	r.writeLock(func() {
		r.AreaId = areaId
	})
}

func (r *Room) GetAreaId() types.Id {
	r.ReadLock()
	defer r.ReadUnlock()
	return r.AreaId
}

func (r *Room) NextLocation(direction types.Direction) types.Coordinate {
	loc := r.GetLocation()
	return loc.Next(direction)
}

func (r *Room) GetExits() []types.Direction {
	r.ReadLock()
	defer r.ReadUnlock()

	exits := make([]types.Direction, len(r.Exits))

	i := 0
	for dir := range r.Exits {
		exits[i] = dir
		i++
	}

	return exits
}

func (r *Room) SetLocked(dir types.Direction, locked bool) {
	r.writeLock(func() {
		if r.HasExit(dir) {
			r.Exits[dir].Locked = locked
		}
	})
}

func (r *Room) IsLocked(dir types.Direction) bool {
	r.ReadLock()
	defer r.ReadUnlock()

	if r.HasExit(dir) {
		return r.Exits[dir].Locked
	}

	return false
}
