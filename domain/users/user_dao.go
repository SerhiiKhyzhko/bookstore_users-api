package users

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/datasources/mysql/users_db"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/date_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUer = "DELETE FROM users WHERE id=&;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreating); getErr != nil{
		return mysqlutils.ParseError(getErr)
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

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreating)
	if saveErr != nil {
		return mysqlutils.ParseError(saveErr)
	}
	
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysqlutils.ParseError(saveErr)
	}
	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id); err != nil {
		return mysqlutils.ParseError(err)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := usersdb.Client.Prepare(queryDeleteUer)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		return mysqlutils.ParseError(err)
	}

	return nil
}