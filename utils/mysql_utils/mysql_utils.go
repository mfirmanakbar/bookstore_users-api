package mysql_utils

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/mfirmanakbar/bookstore_users-api/logger"
	"github.com/mfirmanakbar/bookstore_utils-go/rest_errors"
	"strings"
)

// #1. errorNoRows --> a string from mysql error about there's not row where id = x
// #2. 1062 --> a code mysql error about there's column with same value or duplicated value

const (
	ErrorNoRows = "no rows in result set" // #1
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("no error matching given id")
		}
		logger.Info(sqlErr.Error())
		return rest_errors.NewInternalServerError("error parsing database response", errors.New("database error"))
	}

	switch sqlErr.Number {
	case 1062: // #2
		logger.Info(sqlErr.Error())
		return rest_errors.NewBadRequestError("invalid data")
	}
	logger.Info(sqlErr.Error())
	return rest_errors.NewInternalServerError("error parsing database request", errors.New("database error"))
}
