package services

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/crypto_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/date_utils"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

var UsersService UserServiceInterface = &usersService{}

type usersService struct{
}

type UserServiceInterface interface {
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	GetUser(int64) (*users.User, *rest_errors.RestErr)
	UpdateUser(users.User) (*users.User, *rest_errors.RestErr)
	PartialUpdateUser(users.User) (*users.User, *rest_errors.RestErr)
	DeleteUser(int64) *rest_errors.RestErr
	Search(status string) (users.Users, *rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, *rest_errors.RestErr)
}

func (s *usersService)CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {
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

func (s *usersService)GetUser(userId int64) (*users.User, *rest_errors.RestErr) {
	result := users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *usersService)UpdateUser(user users.User) (*users.User, *rest_errors.RestErr) {
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

func (s *usersService)PartialUpdateUser(user users.User) (*users.User, *rest_errors.RestErr) {
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

func (s *usersService)DeleteUser(userId int64) *rest_errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService)Search(status string) (users.Users, *rest_errors.RestErr) {
	dao := &users.User{}
	users, err := dao.FindByStatus(status)

	if err != nil {
		return nil, err
	}
	return  users, nil
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, *rest_errors.RestErr) {
	dao := &users.User{
		Email: request.Email,
		Password: request.Password,
	}

	if err := dao.FindByEmail(); err != nil {
		return nil, err
	}

	if err := cryptoutils.ComparePassword(dao.Password, request.Password); err != nil {
		return nil, err
	}
	
	return dao, nil
}