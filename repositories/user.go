package repositories

import (
	"errors"
	"movie-app-go/entities"
)

type UserRepository struct {
	data []entities.User
}
type UserRepositoryInterface interface {
	Create(user entities.User) error
	GetByUsername(username string) (entities.User, error)
	// Read(id string) (entities.User, error)
	// ReadAll() ([]entities.User, error)
	Update(user entities.User) error
	// Delete(user *entities.User) error

	// UpdateBalance(id string, amount int) error

	// ReadTransaction(id string) error
	// DeleteTransaction(id string) error
}

func NewUserRepository(data []entities.User) UserRepositoryInterface {
	return &UserRepository{
		data: data,
	}
}
func (repo *UserRepository) Create(user entities.User) error {
	for _, existingUser := range repo.data {
		if existingUser.Username == user.Username {
			return errors.New("user with the same Username already exists")
		}
	}
	repo.data = append(repo.data, user)
	return nil
}
func (repo *UserRepository) GetByUsername(username string) (entities.User, error) {
	for _, user := range repo.data {
		if user.Username == username {
			return user, nil
		}
	}
	return entities.User{}, errors.New("EMPTY_DATA")
}

func (repo *UserRepository) Update(user entities.User) error {
	for i, existingUser := range repo.data {
		if existingUser.Username == user.Username {
			repo.data[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}
