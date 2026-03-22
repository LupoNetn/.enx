package auth

type RegisterInput struct {
	Email    string
	Name     string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	User   any        `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}