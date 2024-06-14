package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable              = "users"
	RefreshSessionsTable    = "refresh_sessions"
	OrganizationsTable      = "organizations"
	UsersOrganizationsTable = "users_organizations"
	DatasourcesTable        = "datasources"
	DatasourceUsersTable    = "datasource_users"
)

type ConnectionData struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func NewPostgresDB(data ConnectionData) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", data.Host, data.Port, data.Username,
			data.Password, data.Name, data.SSLMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
