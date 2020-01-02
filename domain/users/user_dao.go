package users

import (
	"fmt"
	"github.com/mfirmanakbar/bookstore_users-api/datasources/mysql/bookstore_users_db"
	"github.com/mfirmanakbar/bookstore_users-api/utils/date_utils"
	"github.com/mfirmanakbar/bookstore_users-api/utils/errors"
	"strings"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, created_at) VALUES (?, ?, ?, ?);"
	queryGetUser     = "SELECT id, first_name, last_name, email, created_at FROM users WHERE id=?"
)

var (
	usersDB = make(map[int64]*User)
)

// #1. the purpose of pointer as * is to make us able to modified this user object
// #2. Defer statements are generally used to ensure that the files are closed when your work is finished with them,
// 	   or to close the channel, or to catch the panics in the program.
// #3. email_UNIQUE --> from http 500, when we save user with same email, then I custom to http 400
func (user *User) Save() *errors.RestErr { // #1
	stmt, err := bookstore_users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close() // #2

	user.CreatedAt = date_utils.GetNowString()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) { // #3
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))
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

// #1. Scan --> for read all columns from const queryGetUser by sequence
//     &user.Id --> the pointer `&` --> it means to take and POPULATE (Mengisi)
//	   the point is the pointer `&` to make us able to modified User as the method already
func (user *User) Get() *errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt); err != nil { // #1
		fmt.Println(err)
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user id: %d %s", user.Id, err.Error()))
	}
	return nil
}
