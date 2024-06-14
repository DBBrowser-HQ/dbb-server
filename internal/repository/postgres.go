package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
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

	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.Now().Add(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			err = db.Ping()
			if err == nil {
				return db, nil
			}
			if time.Now().After(timeout) {
				return nil, errors.New("timed out waiting for connection: " + err.Error())
			}
		}
	}
}
