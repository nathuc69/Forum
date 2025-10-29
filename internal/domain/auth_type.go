package domain

type GitHubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type AuthRepository interface {
	LoginAuth(username string) error
	RegisterAuth(username, email string) error
	UserExisting(username, email string) bool
	LoginAuthByUsername(Token, username string) error
	LoginAuthByEmail(email string) error
	TokenByEmail(Token, email string) error
}

type AuthService interface {
	GitHub(username, email string) error
	AuthToken(Token, username, email string) error
	Google(name, email string) error
}
