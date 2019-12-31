package services

import (
	"github.com/mfirmanakbar/bookstore_users-api/domain/users"
	"github.com/mfirmanakbar/bookstore_users-api/utils/erros"
)

func CreateUser(user users.User) (*users.User, *erros.RestErr) {

	return &user, nil

	/*return nil, &erros.RestErr{
		Message: "",
		Status:  http.StatusInternalServerError,
		Error:   "",
	}*/
}
