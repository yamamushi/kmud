package types

type Account struct {
	Username             string   `json:"username,omitempty"`
	Email                string   `json:"email,omitempty"`
	HashedPass           string   `json:"hashedpass,omitempty"`
	Groups               []string `json:"groups,omitempty"`
	Permissions          []string `json:"permissions,omitempty"`
	Characters           []string `json:"characters,omitempty"`
	Locked               string   `json:"locked,omitempty"`
	Token                string   `json:"token,omitempty"`
	RequirePasswordReset string   `json:"requirepasswordreset,omitempty"`
}
