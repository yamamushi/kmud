package types

type Account struct {
	Username    string
	Email       string
	HashedPass  string
	Permissions []string
	Characters  []string
}
