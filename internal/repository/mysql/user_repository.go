package mysql

import "database/sql"

type mysqlUserRepository struct {
	*sql.DB
}

func (m *mysqlUserRepository) Create() {}

func (m *mysqlUserRepository) GetByEmail(email string) {}