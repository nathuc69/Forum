package repositories

import (
	"database/sql"
	"forum/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// créer nouvel utilisateur:
func (r *userRepository) Create(user *domain.User) error {
	res, err := r.db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	// Récupérer l’ID auto-généré
	id, err := res.LastInsertId()
	if err == nil {
		user.ID = id
	}
	return nil
}

// récupérer un utilisateur de la BdD par son ID:
func (r *userRepository) GetUserByID(id int) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", id)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// récupérer un utilisateur de la BdD par son email
func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *userRepository) GetByUsername(Logusername string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", Logusername)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) InsertToken(Token, email string) error {
	_, err := r.db.Exec("UPDATE users SET token = ? WHERE email = ? ;", Token, email)
	return err
}

func (r *userRepository) GetUserByToken(Token string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username FROM users WHERE Token = ?", Token)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username /*, &user.Email, &user.Password*/)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteTokenLog(Token string) error {
	_, err := r.db.Exec("UPDATE users SET token = NULL WHERE Token = ?", Token)
	if err != nil {
		return nil
	}
	return err
}
