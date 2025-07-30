package services

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/crypto_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/date_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

var UsersService UserServiceInterface = &usersService{}

type usersService struct{
}

type UserServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	UpdateUser(users.User) (*users.User, *errors.RestErr)
	PartialUpdateUser(users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	Search(status string) (users.Users, *errors.RestErr)
}

func (s *usersService)CreateUser(user users.User) (*users.User, *errors.RestErr) {
	user.Status = users.StatusActive
	user.DateCreating = dateutils.GetNowDbFormat()
	hashedPassword, err := cryptoutils.GetBcrypt(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := user.Validate(); err != nil{
		return nil, err
	}

	if err = user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *usersService)GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *usersService)UpdateUser(user users.User) (*users.User, *errors.RestErr) {
	current, err := s.GetUser(user.Id) 
	if err != nil {
		return nil, err
	}

	current.FirstName = user.FirstName
	current.LastName = user.LastName
	current.Email = user.Email
	

	if err = current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService)PartialUpdateUser(user users.User) (*users.User, *errors.RestErr) {
	current, err := s.GetUser(user.Id) 
	if err != nil {
		return nil, err
	}

	if user.FirstName != "" {
		current.FirstName = user.FirstName
	}
	if user.LastName != "" {
		current.LastName = user.LastName
	}
	if user.Email != "" {
		current.Email = user.Email
	}

	if err = current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService)DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService)Search(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	users, err := dao.FindByStatus(status)

	if err != nil {
		return nil, err
	}
	return  users, nil
}