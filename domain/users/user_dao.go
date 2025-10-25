package users

import (
	"fmt"
	"strings"
	"errors"

	"github.com/SerhiiKhyzhko/bookstore_users-api/datasources/mysql/users_db"
	"github.com/SerhiiKhyzhko/bookstore_users-api/logger"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/mysql_utils"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser    = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser      = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser   = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUer    = "DELETE FROM users WHERE id=&;"
	queryUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmail = "SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE email=? AND status=?;"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare GET user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status); getErr != nil{
		
		logger.Error("error when trying to GET user by id", getErr)	
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	
	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare Save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(
		user.FirstName, user.LastName, user.Email, user.DateCreating, user.Status, user.Password)
	if saveErr != nil {
		
		logger.Error("error when trying to Save user", saveErr)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	
	userId, err := insertResult.LastInsertId()
	if err != nil {

		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	user.Id = userId

	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare Update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id); err != nil {
		
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryDeleteUer)
	if err != nil {
		logger.Error("error when trying to prepare Delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {

		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := usersdb.Client.Prepare(queryUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user by status", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
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

func (user *User) FindByEmail() *rest_errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryFindByEmail)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to find user by email", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, StatusActive)
	if getErr := result.Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating, &user.Status, &user.Password); getErr != nil{
		if strings.Contains(getErr.Error(), mysqlutils.ErrorNoRows) {
			fmt.Println("getErr TEXT", getErr.Error())
			return rest_errors.NewNotFoundError("invalid user credentials")
		}

		logger.Error("error when trying to get user by email and password", getErr)	
		return rest_errors.NewInternalServerError("error when trying to find user by email", errors.New("database error"))
	}
	
	return nil
}