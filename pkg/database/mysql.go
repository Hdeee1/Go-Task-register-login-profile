package database

import (
	"database/sql"
)

func ConnectMySQL() *sql.DB {
	conn, err := sql.Open("mysql", "sihad:dbmysql@/db_register_login_profile")
	if err != nil {
		panic("Failed to connect to database, error" + err.Error())
	}

	return conn
}