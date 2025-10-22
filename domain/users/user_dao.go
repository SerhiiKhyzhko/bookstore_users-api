package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/logger"
	mysqlutils "github.com/SerhiiKhyzhko/bookstore_users-api/utils/mysql_utils"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser        = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser        = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryPartialUpdateUser = "UPDATE users SET"
	queryDeleteUer         = "DELETE FROM users WHERE id=&;"
	queryUserByStatus      = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmail       = "SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE email=? AND status=?;"
)

type UserDaoInterface interface {
	Get(context.Context, int64) (*User, *rest_errors.RestErr)
	Save(context.Context, User) (int64, *rest_errors.RestErr)
	Delete(context.Context, int64) *rest_errors.RestErr
	FindByStatus(context.Context, string) ([]User, *rest_errors.RestErr)
	FindByEmail(context.Context, string) (*User, *rest_errors.RestErr)
	PartialUpdate(context.Context, PartialUser) *rest_errors.RestErr
	Update(context.Context, User) *rest_errors.RestErr
}

type userDaoStruct struct {
	client *sql.DB
}

func NewUserDao(dbClient *sql.DB) *userDaoStruct {
	return &userDaoStruct{
		client: dbClient,
	}
}

func (d *userDaoStruct) Get(ctx context.Context, id int64) (*User, *rest_errors.RestErr) {
	var user User

	result := d.client.QueryRowContext(ctx, queryGetUser, id)
	if getErr := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status); getErr != nil {

		logger.Error("error when trying to GET user by id", getErr)
		return nil, rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}

	return &user, nil
}

func (d *userDaoStruct) Save(ctx context.Context, user User) (int64, *rest_errors.RestErr) {
	insertResult, saveErr := d.client.ExecContext(
		ctx,
		queryInsertUser,
		user.FirstName, user.LastName, user.Email, user.DateCreating, user.Status, user.Password,
	)

	if saveErr != nil {
		logger.Error("error when trying to Save user", saveErr)
		return 0, rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {

		logger.Error("error when trying to get last insert id after creating a new user", err)
		return 0, rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	return userId, nil
}

func (d *userDaoStruct) Update(ctx context.Context, user User) *rest_errors.RestErr {
	if _, err := d.client.ExecContext(
		ctx,
		queryUpdateUser,
		user.FirstName, user.LastName, user.Email, user.Id,
	); err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}

	return nil
}

func (d *userDaoStruct) PartialUpdate(ctx context.Context, user PartialUser) *rest_errors.RestErr {
	query, queryArgs := queryBuilder(queryPartialUpdateUser, user)
	if _, err := d.client.ExecContext(
		ctx,
		query,
		queryArgs...,
	); err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}

	return nil
}

func (d *userDaoStruct) Delete(ctx context.Context, id int64) *rest_errors.RestErr {
	if _, err := d.client.ExecContext(
		ctx,
		queryDeleteUer,
		id,
	); err != nil {

		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}

	return nil
}

func (d *userDaoStruct) FindByStatus(ctx context.Context, status string) ([]User, *rest_errors.RestErr) {
	rows, err := d.client.QueryContext(ctx, queryUserByStatus, status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user by status", errors.New("database error"))
	}
	defer rows.Close()

	result := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find user by status", errors.New("database error"))
		}

		result = append(result, user)
	}
	if len(result) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no user matching status %v", status))
	}
	return result, nil
}

func (d *userDaoStruct) FindByEmail(ctx context.Context, email string) (*User, *rest_errors.RestErr) {
	var user User
	result := d.client.QueryRowContext(
		ctx,
		queryFindByEmail,
		email, StatusActive,
	)
	if getErr := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status, &user.Password); getErr != nil {
		if strings.Contains(getErr.Error(), mysqlutils.ErrorNoRows) {
			return errors.NewNotFoundError("invalid user credentials")
		}

		logger.Error("error when trying to get user by email and password", getErr)
		return nil, rest_errors.NewInternalServerError("error when trying to find user by email", errors.New("database error"))
	}

	return &user, nil
}

func queryBuilder(rawQuery string, userData PartialUser) (string, []any) {
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
	data = append(data, userData.Id)
	return fmt.Sprintf("%s %s %s", rawQuery, strings.Join(queryPart, ", "), setLast), data
}
