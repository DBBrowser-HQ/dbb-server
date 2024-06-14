package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
)

type DataSourcePostgres struct {
	db *sqlx.DB
}

func NewDataSourcePostgres(db *sqlx.DB) *DataSourcePostgres {
	return &DataSourcePostgres{db: db}
}

func (r *DataSourcePostgres) GetUser(userId int) (string, string, error) {
	query := fmt.Sprintf(`SELECT login, password_hash FROM %s WHERE id = $1`, UsersTable)

	var input struct {
		Login        string `db:"login"`
		PasswordHash string `db:"password_hash"`
	}
	err := r.db.Get(&input, query, userId)
	if err != nil {
		return "", "", err
	}
	return input.Login, input.PasswordHash, nil
}

func (r *DataSourcePostgres) GetUnusedPort() (int, error) {
	query := fmt.Sprintf(`SELECT port + 1 FROM %s
				ORDER BY port DESC
				LIMIT 1`, DatasourcesTable)

	var port int
	err := r.db.Get(&port, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			port, err = strconv.Atoi(os.Getenv("DB_PORT")) // default port value
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return port, err
}

func (r *DataSourcePostgres) CheckDatasourceExistence(dbName string, organizationId int) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE name=$1 AND organization_id=$2`, DatasourcesTable)
	var count int
	err := r.db.Get(&count, query, dbName, organizationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}
	if count == 0 {
		return false, nil
	}
	return true, err
}

func (r *DataSourcePostgres) CreateDatasource(dbHost, dbName string, dbPort, organizationId int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (host, port, name, organization_id)
								VALUES ($1, $2, $3, $4) RETURNING id`, DatasourcesTable)

	var id int
	err := r.db.Get(&id, query, dbHost, dbPort, dbName, organizationId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DataSourcePostgres) AddRolesData(usernames, passwords []string, datasourceId int) error {
	query := fmt.Sprintf(`INSERT INTO %s (username, password, datasource_id) VALUES ($1, $2, $3)`, DatasourceUsersTable)

	if len(usernames) != len(passwords) {
		return errors.New("usernames and passwords have different length")
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	for i := range usernames {
		_, err = stmt.Exec(usernames[i], passwords[i], datasourceId)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *DataSourcePostgres) CreateRolesInH2Db(h2Db *sqlx.DB) error {
	queryCreateReaderRole := `CREATE ROLE READER`
	queryCreateRedactorRole := `CREATE ROLE REDACTOR`
	queryGrantRightsReader := `GRANT SELECT ON SCHEMA PUBLIC TO READER`
	queryGrantRightsRedactor := `GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA PUBLIC TO REDACTOR`

	defer h2Db.Close()

	tx, err := h2Db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(queryCreateReaderRole); err != nil {
		return err
	}
	if _, err = tx.Exec(queryGrantRightsReader); err != nil {
		return err
	}
	if _, err = tx.Exec(queryCreateRedactorRole); err != nil {
		return err
	}
	if _, err = tx.Exec(queryGrantRightsRedactor); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DataSourcePostgres) GetHostNames() ([]string, error) {
	query := fmt.Sprintf(`SELECT host FROM %s`, DatasourcesTable)
	var hosts []string
	err := r.db.Select(&hosts, query)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}
