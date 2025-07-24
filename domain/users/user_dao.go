package users

import (
	"fmt"
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/datasources/mysql/users_db"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/date_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows = "no rows in result set"
	queryInertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err = result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating); err != nil{
		if strings.Contains(err.Error(), errorNoRows){
			return errors.NewNotFoundError(fmt.Sprintf("user %v not found", user.Id))
		}
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to get user %v, %s", user.Id, err.Error()))
	}
	
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryInertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreating = dateutils.GetNowString()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreating)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail){
			return errors.NewBadRequestError(fmt.Sprintf("email %v already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}
	user.Id = userId

	return nil
}