package types

type AuthRequest struct {
	Secret     string `json:"secret"`
	Username   string `json:"username"`
	HashedPass string `json:"hashedpass"`
}

type AuthResponse struct {
	AuthToken string `json:"authtoken"`
	Err       string `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}

type AccountInfoRequest struct {
	Secret string `json:"secret"`
	Token  string `json:"token"`
	Field  string `json:"field"`
}

type AccountInfoResponse struct {
	Account Account `json:"account"`
	Err     string  `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}

type AccountRegistrationRequest struct {
	Secret     string `json:"secret"`
	Username   string `json:"username"`
	HashedPass string `json:"hashedpass"`
	Email      string `json:"email"`
}

type AccountRegistrationResponse struct {
	Err string `json:"error"` // errors don't JSON-marshal, so we use a string
}

type SearchRequest struct {
	Secret  string  `json:"secret"`
	Token   string  `json:"token"`
	Account Account `json:"account"`
}

type SearchResponse struct {
	Accounts []Account `json:"accounts"`
	Err      string    `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}

type ModifyRequest struct {
	Secret  string  `json:"secret"`
	Token   string  `json:"token"`
	Account Account `json:"account"`
}

type ModifyResponse struct {
	Account Account `json:"account"`
	Err     string  `json:"error,omitempty"` // errors don't JSON-marshal, so we use a string
}
