package services

import (
	"context"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/domain/users"
	cryptoutils "github.com/SerhiiKhyzhko/bookstore_users-api/v2/utils/crypto_utils"
	dateutils "github.com/SerhiiKhyzhko/bookstore_users-api/v2/utils/date_utils"
)

type usersService struct {
	userDao    users.UserDaoInterface
	ctxTimeout time.Duration
}

type UserServiceInterface interface {
	CreateUser(context.Context, users.User) (*users.User, error)
	GetUser(context.Context, int64) (*users.User, error)
	UpdateUser(context.Context, users.User) (*users.User, error)
	PartialUpdateUser(context.Context, users.PartialUser) (*users.User, error)
	DeleteUser(context.Context, int64) error
	Search(context.Context, string) (users.Users, error)
	LoginUser(context.Context, users.LoginRequest) (*users.User, error)
}

func NewUserService(timeout time.Duration, dao users.UserDaoInterface) *usersService {
	return &usersService{
		userDao:    dao,
		ctxTimeout: timeout,
	}
}

func (s *usersService) CreateUser(ctx context.Context, user users.User) (*users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()
	user.Status = users.StatusActive
	user.DateCreating = dateutils.GetNowDbFormat()
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := cryptoutils.GetBcrypt(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	id, err := s.userDao.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = id

	return &user, nil
}

func (s *usersService) GetUser(ctx context.Context, userId int64) (*users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()
	result, err := s.userDao.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) UpdateUser(ctx context.Context, user users.User) (*users.User, error) {
	updateCtx, cancelUpdate := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancelUpdate()
	if err := s.userDao.Update(updateCtx, user); err != nil {
		return nil, err
	}

	getCtx, cancelGet := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancelGet()
	return s.GetUser(getCtx, user.Id)
}

func (s *usersService) PartialUpdateUser(ctx context.Context, user users.PartialUser) (*users.User, error) {
	updateCtx, cancelUpdate := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancelUpdate()
	if err := s.userDao.PartialUpdate(updateCtx, user); err != nil {
		return nil, err
	}

	getCtx, cancelGet := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancelGet()
	return s.GetUser(getCtx, user.Id)
}

func (s *usersService) DeleteUser(ctx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()
	return s.userDao.Delete(ctx, userId)
}

func (s *usersService) Search(ctx context.Context, status string) (users.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()
	users, err := s.userDao.FindByStatus(ctx, status)

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *usersService) LoginUser(ctx context.Context, request users.LoginRequest) (*users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()
	user, err := s.userDao.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if err := cryptoutils.ComparePassword(user.Password, request.Password); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}
