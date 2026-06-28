package handler

import (
	"strings"

	"github.com/go-sql-driver/mysql"
)

func isDupErr(err error) bool {
	if me, ok := err.(*mysql.MySQLError); ok {
		return me.Number == 1062
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
