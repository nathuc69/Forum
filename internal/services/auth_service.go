package services

import (
	"fmt"
	"forum/internal/domain"
)

type authService struct {
	repo domain.AuthRepository
}

func NewAuthService(repo domain.AuthRepository) domain.AuthService {
	return &authService{repo: repo}
}

func (s *authService) GitHub(username, email string) error {
	Exist := s.repo.UserExisting(username, email)
	if Exist {
		err := s.repo.LoginAuth(username)
		if err != nil {
			fmt.Println("erreur login ")
			return err
		}
		return nil
	}
	err := s.repo.RegisterAuth(username, email)
	if err != nil {
		fmt.Println("erreur register")
		return err
	}
	return nil
}

func (s *authService) AuthToken(Token, username, email string) error {
	if email == "" {
		err := s.repo.LoginAuthByUsername(Token, username)
		if err != nil {
			return err
		}
	} else {
		err := s.repo.TokenByEmail(Token, email)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *authService) Google(name, email string) error {
	if email == "" {
		email = "test"
	}
	Exist := s.repo.UserExisting(name, email)
	if Exist {
		err := s.repo.LoginAuthByEmail(email)
		if err != nil {
			fmt.Println("erreur login ")
			return err
		}
		return nil
	}
	err := s.repo.RegisterAuth(name, email)
	if err != nil {
		fmt.Println("erreur register")
		return err
	}
	return nil
}
