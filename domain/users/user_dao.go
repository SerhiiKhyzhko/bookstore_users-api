package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/user_errors"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
)

const (
	queryInsertUser           = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser              = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser           = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	rawQueryPartialUpdateUser = "UPDATE users SET"
	queryDeleteUser           = "DELETE FROM users WHERE id=?;"
	queryUserByStatus         = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmail          = "SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE email=? AND status=?;"
)

type UserDaoInterface interface {
	Get(context.Context, int64) (*User, error)
	Save(context.Context, User) (int64, error)
	Delete(context.Context, int64) error
	FindByStatus(context.Context, string) ([]User, error)
	FindByEmail(context.Context, string) (*User, error)
	PartialUpdate(context.Context, PartialUser) error
	Update(context.Context, User) error
}

type userDaoStruct struct {
	client *sql.DB
	logger *logger.Logger
}

func NewUserDao(dbClient *sql.DB, logger *logger.Logger) *userDaoStruct {
	return &userDaoStruct{
		client: dbClient,
		logger: logger,
	}
}

func (d *userDaoStruct) Get(ctx context.Context, id int64) (*User, error) {
	var user User

	result := d.client.QueryRowContext(ctx, queryGetUser, id)
	if getErr := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status); getErr != nil {

		if errors.Is(getErr, sql.ErrNoRows) {
			return nil, user_errors.NotFoundErr
		}
		if errors.Is(getErr, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return nil, user_errors.RequestTimeoutErr
		}

		d.logger.Error("error when trying to GET user by id", getErr)
		return nil, fmt.Errorf("error when trying to get user: %w", getErr)
	}

	return &user, nil
}

func (d *userDaoStruct) Save(ctx context.Context, user User) (int64, error) {
	insertResult, saveErr := d.client.ExecContext(
		ctx,
		queryInsertUser,
		user.FirstName, user.LastName, user.Email, user.DateCreating, user.Status, user.Password,
	)

	if saveErr != nil {
		if errors.Is(saveErr, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return 0, user_errors.RequestTimeoutErr
		}
		d.logger.Error("error when trying to Save user", saveErr)
		return 0, fmt.Errorf("save failed %w", saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {

		d.logger.Error("error when trying to get last insert id after creating a new user", err)
		return 0, fmt.Errorf("error when trying to save user: %w", err)
	}
	return userId, nil
}

func (d *userDaoStruct) Update(ctx context.Context, user User) error {
	if _, err := d.client.ExecContext(
		ctx,
		queryUpdateUser,
		user.FirstName, user.LastName, user.Email, user.Id,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return user_errors.RequestTimeoutErr
		}
		d.logger.Error("error when trying to update user", err)
		return fmt.Errorf("error when trying to update user: %w", err)
	}

	return nil
}

func (d *userDaoStruct) PartialUpdate(ctx context.Context, user PartialUser) error {
	query, queryArgs, queryErr := queryBuilder(rawQueryPartialUpdateUser, user)
	if queryErr != nil {
		return queryErr
	}
	if _, err := d.client.ExecContext(
		ctx,
		query,
		queryArgs...,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return user_errors.RequestTimeoutErr
		}
		d.logger.Error("error when trying to update user", err)
		return fmt.Errorf("error when trying to update user: %w", err)
	}

	return nil
}

func (d *userDaoStruct) Delete(ctx context.Context, id int64) error {
	if _, err := d.client.ExecContext(
		ctx,
		queryDeleteUser,
		id,
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return user_errors.RequestTimeoutErr
		}
		d.logger.Error("error when trying to delete user", err)
		return fmt.Errorf("error when trying to delete user:%w", err)
	}

	return nil
}

func (d *userDaoStruct) FindByStatus(ctx context.Context, status string) ([]User, error) {
	rows, err := d.client.QueryContext(ctx, queryUserByStatus, status)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return nil, user_errors.RequestTimeoutErr
		}
		d.logger.Error("error when trying to find users by status", err)
		return nil, fmt.Errorf("error when trying to find user by status: %w", err)
	}
	defer rows.Close()

	result := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status); err != nil {
			d.logger.Error("error when scan user row into user struct", err)
			return nil, fmt.Errorf("error when trying to find user by status: %w", err)
		}

		result = append(result, user)
	}

	if err = rows.Err(); err != nil {
		d.logger.Error("error during rows iteration for FindByStatus", err)
		return nil, fmt.Errorf("error processing user list: %w", err)
	}

	if len(result) == 0 {
		return nil, user_errors.NotFoundErr
	}
	return result, nil
}

func (d *userDaoStruct) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := d.client.QueryRowContext(
		ctx,
		queryFindByEmail,
		email, StatusActive,
	)
	if getErr := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status, &user.Password); getErr != nil {
		if errors.Is(getErr, sql.ErrNoRows) {
			return nil, user_errors.NotFoundErr //rest_errors.NewNotFoundError("invalid user credentials")
		}
		if errors.Is(getErr, context.DeadlineExceeded) {
			d.logger.Error("db request timeout", context.DeadlineExceeded)
			return nil, user_errors.RequestTimeoutErr
		}

		d.logger.Error("error when trying to get user by email and password", getErr)
		return nil, fmt.Errorf("error when trying to find user by email: %w", getErr)
	}

	return &user, nil
}

func queryBuilder(rawQuery string, userData PartialUser) (string, []any, error) {
	const (
		setFirstName = "first_name = ?"
		setLastName  = "last_name = ?"
		setEmail     = "email = ?"
		setStatus    = "status = ?"
		setLast      = "WHERE id=?;"
	)

	queryPart := make([]string, 0, 4)
	data := make([]any, 0, 5)
	if userData.FirstName != nil {
		data = append(data, *userData.FirstName)
		queryPart = append(queryPart, setFirstName)
	}
	if userData.LastName != nil {
		data = append(data, *userData.LastName)
		queryPart = append(queryPart, setLastName)
	}
	if userData.Email != nil {
		data = append(data, *userData.Email)
		queryPart = append(queryPart, setEmail)
	}
	if userData.Status != nil {
		data = append(data, *userData.Status)
		queryPart = append(queryPart, setStatus)
	}
	if len(data) == 0 {
		return "", data, user_errors.BadRequestErr
	}
	data = append(data, userData.Id)
	return fmt.Sprintf("%s %s %s", rawQuery, strings.Join(queryPart, ", "), setLast), data, nil
}
