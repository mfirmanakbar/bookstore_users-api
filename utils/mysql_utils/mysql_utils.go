package mysql_utils

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/mfirmanakbar/bookstore_users-api/utils/errors"
	"strings"
)

// #1. errorNoRows --> a string from mysql error about there's not row where id = x
// #2. 1062 --> a code mysql error about there's column with same value or duplicated value

const (
	errorNoRows = "no rows in result set" // #1
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no error matching given id")
		}
		return errors.NewInternalServerError(fmt.Sprintf("error parsing database response from %s", sqlErr.Error()))
	}

	switch sqlErr.Number {
	case 1062: // #2
		return errors.NewBadRequestError(fmt.Sprintf("invalid data from %s", sqlErr.Error()))
	}
	return errors.NewInternalServerError(fmt.Sprintf("error parsing database request from %s", sqlErr.Error()))
}
