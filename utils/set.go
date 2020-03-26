package utils

type Set map[string]bool

func (s Set) Contains(key string) bool {
	_, found := s[key]
	return found
}

func (s Set) Insert(key string) {
	s[key] = true
}

func (s Set) Remove(key string) {
	delete(s, key)
}

func (s Set) Size() int {
	return len(s)
}
