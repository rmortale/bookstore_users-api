package mysql_utils

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/rmortale/bookstore_oauth-go/oauth/errors"
	"github.com/rmortale/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	fmt.Println(err)
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		return rest_errors.NewInternalServerError("error parsing database response", errors.NewError("error parsing database response"))
	}
	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("invalid data")
	}
	return rest_errors.NewInternalServerError("error processing request", errors.NewError("error parsing database response"))
}
