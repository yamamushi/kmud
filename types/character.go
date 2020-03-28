package types

type Character interface {
	Object
	Nameable
	Locateable
	Container
	SetRoomId(Id)
	Hit(int)
	Heal(int)
	GetHitPoints() int
	SetHitPoints(int)
	GetHealth() int
	SetHealth(int)
	GetSkills() []Id
	AddSkill(Id)
}

type CharacterList []Character

func (l CharacterList) Names() []string {
	names := make([]string, len(l))
	for i, char := range l {
		names[i] = char.GetName()
	}
	return names
}