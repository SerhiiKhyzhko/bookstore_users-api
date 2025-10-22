package services

import (
	"context"

	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/crypto_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/date_utils"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

type usersService struct {
	userDao users.UserDaoInterface
}

type UserServiceInterface interface {
	CreateUser(context.Context, users.User) (*users.User, *rest_errors.RestErr)
	GetUser(context.Context, int64) (*users.User, *rest_errors.RestErr)
	UpdateUser(context.Context, users.User) *rest_errors.RestErr
	PartialUpdateUser(context.Context, users.User) *rest_errors.RestErr
	DeleteUser(context.Context, int64) *rest_errors.RestErr
	Search(context.Context, string) (users.Users, *rest_errors.RestErr)
	LoginUser(context.Context, users.LoginRequest) (*users.User, *rest_errors.RestErr)
}

func NewUserService(dao users.UserDaoInterface) *usersService {
	return &usersService{
		userDao: dao,
	}
}

func (s *usersService) CreateUser(ctx context.Context, user users.User) (*users.User, *rest_errors.RestErr) {
	user.Status = users.StatusActive
	user.DateCreating = dateutils.GetNowDbFormat()
	hashedPassword, err := cryptoutils.GetBcrypt(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := user.Validate(); err != nil {
		return nil, err
	}

	id, err := s.userDao.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = id

	return &user, nil
}

func (s *usersService) GetUser(ctx context.Context, userId int64) (*users.User, *rest_errors.RestErr) {
	result, err := s.userDao.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) UpdateUser(ctx context.Context, user users.User)  *rest_errors.RestErr {
	if err := s.userDao.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *usersService) PartialUpdateUser(ctx context.Context, user users.PartialUser) *rest_errors.RestErr {
	if err := s.userDao.PartialUpdate(ctx, user); err != nil {
		return err
	}
	
	return nil
}

func (s *usersService) DeleteUser(ctx context.Context, userId int64) *rest_errors.RestErr {
	return s.userDao.Delete(ctx, userId)
}

func (s *usersService) Search(ctx context.Context, status string) (users.Users, *rest_errors.RestErr) {
	users, err := s.userDao.FindByStatus(ctx, status)

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *usersService) LoginUser(ctx context.Context, request users.LoginRequest) (*users.User, *rest_errors.RestErr) {
	user, err := s.userDao.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if err := cryptoutils.ComparePassword(user.Password, request.Password); err != nil {
		return nil, err
	}

	return user, nil
}
