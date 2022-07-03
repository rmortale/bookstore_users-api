package users

import (
	"fmt"
	"github.com/rmortale/bookstore_users-api/datasources/mysql/users_db"
	"github.com/rmortale/bookstore_users-api/utils/errors"
	"github.com/rmortale/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser       = "insert into users(first_name, last_name, email, date_created, status, password) values(?,?,?,?,?,?);"
	queryGetUser          = "select id,first_name, last_name, email, date_created, status from users where id=?;"
	queryUpdateUser       = "update users set first_name=?, last_name=?, email=? where id=?;"
	queryDeleteUser       = "delete from users where id=?;"
	queryFindUserByStatus = "select id,first_name, last_name, email, date_created, status from users where status=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	user.Id = id
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no user matching status %s", status))
	}
	return results, nil
}
