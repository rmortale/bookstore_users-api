package users

import (
	"errors"
	"fmt"
	"github.com/rmortale/bookstore_utils-go/rest_errors"

	"github.com/rmortale/bookstore_users-api/datasources/mysql/users_db"
	"github.com/rmortale/bookstore_users-api/logger"
	"github.com/rmortale/bookstore_users-api/utils/mysql_utils"
	"strings"
)

const (
	queryInsertUser             = "insert into users(first_name, last_name, email, date_created, status, password) values(?,?,?,?,?,?);"
	queryGetUser                = "select id,first_name, last_name, email, date_created, status from users where id=?;"
	queryUpdateUser             = "update users set first_name=?, last_name=?, email=? where id=?;"
	queryDeleteUser             = "delete from users where id=?;"
	queryFindUserByStatus       = "select id,first_name, last_name, email, date_created, status from users where status=?;"
	queryFindByEmailAndPassword = "select id,first_name, last_name, email, date_created, status from users where email=? and password=? and status=?;"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare user statement", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when getting last insert id", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	user.Id = id
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when preparing client for update", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when updating user", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when preparing client for delete", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when deleting user", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when preparing client for search user", err)
		return nil, rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when searching user", err)
		return nil, rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scanning results", err)
			return nil, rest_errors.NewInternalServerError("database error", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no user matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		if strings.Contains(err.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("no user found with given credentials")
		}
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.NewInternalServerError("database error", errors.New("database error"))
	}
	return nil
}
