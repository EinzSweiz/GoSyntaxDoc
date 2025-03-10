package user

import (
	"GoSyntaxDoc/domain/entities"
	"GoSyntaxDoc/infrastructure/repositories"
	"fmt"

	"github.com/sirupsen/logrus"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) FetchUserById(userID int) (*entities.User, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := s.Repo.FetchUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil // ✅ Return the user object instead of an HTTP response
}

func (s *UserService) HandleUserCreated(firstName string, lastName string) (*entities.User, error) {
	if firstName == "" || lastName == "" {
		logrus.WithFields(logrus.Fields{"error": "Invalid user data"}).Error("Invalid user data")
		return nil, fmt.Errorf("first name and last name are required")
	}

	user, err := s.Repo.CreateUser(firstName, lastName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Failed to create user")
		return nil, err
	}

	return user, nil // ✅ Correctly returning both user and error
}

func (s *UserService) HandleUserRead() ([]entities.User, error) {
	users, err := s.Repo.FetchAllUsers()
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Failed to fetch all users")
		return nil, err
	}
	return users, nil
}
