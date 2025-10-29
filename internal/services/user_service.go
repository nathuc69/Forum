package services

import (
	"errors"
	"fmt"
	"forum/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(username, email, password string) error {
	existing, err := s.repo.GetUserByEmail(email)
	if err == nil && existing != nil {
		return errors.New("❌ email already registered")
	}
	existing2, err2 := s.repo.GetByUsername(username)
	if err2 == nil && existing2 != nil {
		return errors.New("❌ username already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("err ici ")
		return err
	}

	user := domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(&user)
}

func (s *userService) Authenticate(email, password string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		fmt.Println("❌ invalid email")
		return nil, errors.New("❌ invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("❌ invalid password")
		return nil, errors.New("❌ invalid password")
	}
	fmt.Println("✅ hashing password successful")

	return user, nil
}

func (s *userService) TokenLogIn(Token, email string) error {
	err := s.repo.InsertToken(Token, email)
	if err != nil {
		return err
	}
	return nil
}
func (s *userService) Home(Token string) (*domain.User, error) {
	user, err := s.repo.GetUserByToken(Token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return user, nil
}

func (s *userService) Logout(Token string) error {
	err := s.repo.DeleteTokenLog(Token)
	if err != nil {
		return err
	}
	return nil
}
