package users

import (
	"github.com/mfirmanakbar/bookstore_users-api/utils/erros"
	"strings"
)

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// this is the function --> call with users.Validate(&user)
/*func Validate(user *User) *erros.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return erros.NewBadRequestError("invalid json body")
	}
	return nil
}*/

// this is the method --> call with user.Validate()
func (user *User) Validate() *erros.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return erros.NewBadRequestError("invalid json body")
	}
	return nil
}
