package mysqlutils

import (
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
)

const(
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing db responce")
	}
	switch sqlErr.Number {
	case 1062:
			return errors.NewBadRequestError("invalid data")
	}
	return errors.NewInternalServerError("error processing request")
}