package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/user_errors"
	cryptoutils "github.com/SerhiiKhyzhko/bookstore_users-api/v2/utils/crypto_utils"
	"github.com/stretchr/testify/assert"
)

type mockUserDao struct {
	returnUser  *users.User
	returnUsers []users.User
	returnId    int64
	returnErr   error
}

func (m *mockUserDao) Get(_ context.Context, _ int64) (*users.User, error) {
	return m.returnUser, m.returnErr
}

func (m *mockUserDao) Save(_ context.Context, _ users.User) (int64, error) {
	return m.returnId, m.returnErr
}

func (m *mockUserDao) Delete(_ context.Context, _ int64) error {
	return m.returnErr
}

func (m *mockUserDao) FindByStatus(_ context.Context, _ string) ([]users.User, error) {
	return m.returnUsers, m.returnErr
}

func (m *mockUserDao) FindByEmail(_ context.Context, _ string) (*users.User, error) {
	return m.returnUser, m.returnErr
}

func (m *mockUserDao) PartialUpdate(_ context.Context, _ users.PartialUser) error {
	return m.returnErr
}

func (m *mockUserDao) Update(_ context.Context, _ users.User) error {
	return m.returnErr
}

func newService(dao *mockUserDao) UserServiceInterface {
	return NewUserService(2*time.Second, dao)
}

func validUser() users.User {
	return users.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
		Password:  "password123",
	}
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		dao := &mockUserDao{returnId: 1}
		service := newService(dao)

		result, err := service.CreateUser(context.Background(), validUser())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.Id)
		assert.Equal(t, users.StatusActive, result.Status)
		assert.NotEmpty(t, result.DateCreating)
		// result.Password should be hassed 
		assert.NotEqual(t, "password123", result.Password)
	})

	t.Run("InvalidUser_EmptyEmail", func(t *testing.T) {
		dao := &mockUserDao{}
		service := newService(dao)

		input := validUser()
		input.Email = ""

		result, err := service.CreateUser(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_errors.BadRequestErr)
	})

	t.Run("InvalidUser_ShortPassword", func(t *testing.T) {
		dao := &mockUserDao{}
		service := newService(dao)

		input := validUser()
		input.Password = "123"

		result, err := service.CreateUser(context.Background(), input)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, user_errors.BadRequestErr)
	})

	t.Run("DaoError", func(t *testing.T) {
		dao := &mockUserDao{returnErr: errors.New("db error")}
		service := newService(dao)

		result, err := service.CreateUser(context.Background(), validUser())

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("ContextTimeout", func(t *testing.T) {
		dao := &mockUserDao{returnErr: user_errors.RequestTimeoutErr}
		service := newService(dao)

		result, err := service.CreateUser(context.Background(), validUser())

		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_errors.RequestTimeoutErr)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := &users.User{Id: 1, Email: "john@test.com"}
		dao := &mockUserDao{returnUser: expected}
		service := newService(dao)

		result, err := service.GetUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.Id)
		assert.Equal(t, "john@test.com", result.Email)
	})

	t.Run("NotFound", func(t *testing.T) {
		dao := &mockUserDao{returnErr: user_errors.NotFoundErr}
		service := newService(dao)

		result, err := service.GetUser(context.Background(), 99)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_errors.NotFoundErr)
	})

	t.Run("ContextTimeout", func(t *testing.T) {
		dao := &mockUserDao{returnErr: user_errors.RequestTimeoutErr}
		service := newService(dao)

		result, err := service.GetUser(context.Background(), 1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_errors.RequestTimeoutErr)
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// generating password hash to make ComparePassword work
		hash, err := cryptoutils.GetBcrypt("password123")
		assert.NoError(t, err)

		dao := &mockUserDao{
			returnUser: &users.User{
				Id:       1,
				Email:    "john@test.com",
				Password: hash,
			},
		}
		service := newService(dao)

		result, err := service.LoginUser(context.Background(), users.LoginRequest{
			Email:    "john@test.com",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.Id)
		// result.Password field should be empty
		assert.Empty(t, result.Password)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		dao := &mockUserDao{returnErr: user_errors.NotFoundErr}
		service := newService(dao)

		result, err := service.LoginUser(context.Background(), users.LoginRequest{
			Email:    "unknown@test.com",
			Password: "password123",
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_errors.NotFoundErr)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		hash, err := cryptoutils.GetBcrypt("password123")
		assert.NoError(t, err)

		dao := &mockUserDao{
			returnUser: &users.User{
				Id:       1,
				Email:    "john@test.com",
				Password: hash,
			},
		}
		service := newService(dao)

		result, err := service.LoginUser(context.Background(), users.LoginRequest{
			Email:    "john@test.com",
			Password: "wrongpassword",
		})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		dao := &mockUserDao{}
		service := newService(dao)

		err := service.DeleteUser(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("DaoError", func(t *testing.T) {
		dao := &mockUserDao{returnErr: errors.New("db error")}
		service := newService(dao)

		err := service.DeleteUser(context.Background(), 1)

		assert.Error(t, err)
	})
}

func TestPartialUpdateUser(t *testing.T) {
	firstName := "Updated"

	t.Run("Success", func(t *testing.T) {
		expected := &users.User{Id: 1, FirstName: "Updated", Email: "john@test.com"}
		dao := &mockUserDao{returnUser: expected}
		service := newService(dao)

		result, err := service.PartialUpdateUser(context.Background(), users.PartialUser{
			Id:        1,
			FirstName: &firstName,
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated", result.FirstName)
	})

	t.Run("NotFound", func(t *testing.T) {
		dao := &mockUserDao{returnErr: user_errors.NotFoundErr}
		service := newService(dao)

		result, err := service.PartialUpdateUser(context.Background(), users.PartialUser{
			Id:        99,
			FirstName: &firstName,
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_errors.NotFoundErr)
	})

	t.Run("EmptyFields_BadRequest", func(t *testing.T) {
		dao := &mockUserDao{returnErr: user_errors.BadRequestErr}
		service := newService(dao)

		result, err := service.PartialUpdateUser(context.Background(), users.PartialUser{
			Id: 1,
		})

		assert.Nil(t, result)
		assert.ErrorIs(t, err, user_errors.BadRequestErr)
	})
}
