package users

import (
	"errors"
	"fmt"
	"github.com/mfirmanakbar/bookstore_users-api/datasources/mysql/bookstore_users_db"
	"github.com/mfirmanakbar/bookstore_users-api/logger"
	"github.com/mfirmanakbar/bookstore_users-api/utils/mysql_utils"
	"github.com/mfirmanakbar/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, created_at, password, status) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE id=?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE email=? and password=? and status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

// #1. the purpose of pointer as * is to make us able to modified this user object
// #2. Defer statements are generally used to ensure that the files are closed when your work is finished with them,
// 	   or to close the channel, or to catch the panics in the program.
func (user *User) Save() rest_errors.RestErr { // #1
	stmt, err := bookstore_users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	defer stmt.Close() // #2

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt, user.Password, user.Status)
	if saveErr != nil {
		//return mysql_utils.ParseError(saveErr)
		logger.Error("error when trying to save user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		//return mysql_utils.ParseError(saveErr)
		logger.Error("error when trying to get Last Insert Id user after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	user.Id = userId
	return nil
}

// #3. Scan --> for read all columns from const queryGetUser by sequence
//     &user.Id --> the pointer `&` --> it means to take and POPULATE (Mengisi)
//	   the point is the pointer `&` to make us able to modified User as the method already
func (user *User) Get() rest_errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil { // #3
		//return mysql_utils.ParseError(getErr)
		logger.Error("error when trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error when trying to get user", errors.New("database error"))
	}
	return nil
}

func (user *User) Update() rest_errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		//return mysql_utils.ParseError(err)
		logger.Error("error when trying to update user by id", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Id); err != nil {
		//return mysql_utils.ParseError(err)
		logger.Error("error when trying to delete user by id", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	stmt, err := bookstore_users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0) // define a map
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			//return nil, mysql_utils.ParseError(err)
			logger.Error("error when scanning user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil {
		//return mysql_utils.ParseError(getErr)
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			fmt.Println(getErr.Error())
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}
	return nil
}
