package domain

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type AuthRepository interface {
	LoginAuth(username string) error
	RegisterAuth(username, email string) error
	UserExisting(username string) bool
	LoginAuthByUsername(Token, username string) error
}

type AuthService interface {
	GitHub(username, email string) error
	AuthToken(Token, username string) error
	Google(name, email string) error
}
