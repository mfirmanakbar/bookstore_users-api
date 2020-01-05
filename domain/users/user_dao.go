package users

import (
	"fmt"
	"github.com/mfirmanakbar/bookstore_users-api/datasources/mysql/bookstore_users_db"
	"github.com/mfirmanakbar/bookstore_users-api/logger"
	"github.com/mfirmanakbar/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser           = "INSERT INTO users(first_name, last_name, email, created_at, password, status) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser              = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE id=?;"
	queryUpdateUser           = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser           = "DELETE FROM users WHERE id=?;"
	queryFindUserUserByStatus = "SELECT id, first_name, last_name, email, created_at, status FROM users WHERE status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

// #1. the purpose of pointer as * is to make us able to modified this user object
// #2. Defer statements are generally used to ensure that the files are closed when your work is finished with them,
// 	   or to close the channel, or to catch the panics in the program.
func (user *User) Save() *errors.RestErr { // #1
	stmt, err := bookstore_users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close() // #2

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.CreatedAt, user.Password, user.Status)
	if saveErr != nil {
		//return mysql_utils.ParseError(saveErr)
		logger.Error("error when trying to save user", err)
		return errors.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		//return mysql_utils.ParseError(saveErr)
		logger.Error("error when trying to get Last Insert Id user after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

// #3. Scan --> for read all columns from const queryGetUser by sequence
//     &user.Id --> the pointer `&` --> it means to take and POPULATE (Mengisi)
//	   the point is the pointer `&` to make us able to modified User as the method already
func (user *User) Get() *errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); getErr != nil { // #3
		//return mysql_utils.ParseError(getErr)
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		//return mysql_utils.ParseError(err)
		logger.Error("error when trying to update user by id", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := bookstore_users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Id); err != nil {
		//return mysql_utils.ParseError(err)
		logger.Error("error when trying to delete user by id", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := bookstore_users_db.Client.Prepare(queryFindUserUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0) // define a map
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.Status); err != nil {
			logger.Error("error when scanning user row into user struct", err)
			return nil, errors.NewInternalServerError("database error")
			//return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching %s", status))
	}
	return results, nil
}
