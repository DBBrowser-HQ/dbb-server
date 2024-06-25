package repository

import (
	"dbb-server/internal/model"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ContainerDB struct {
	db *sqlx.DB
}

func NewContainerDB(dbHost, dbUsername, dbPassword, dbName string, dbPort int) (*ContainerDB, error) {
	connectionString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		dbHost, dbPort, dbUsername, dbPassword, dbName, "disable")

	conn, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.Now().Add(15 * time.Second)

	for {
		select {
		case <-ticker.C:
			err = conn.Ping()
			if err == nil {
				return &ContainerDB{conn}, nil
			}
			if time.Now().After(timeout) {
				return nil, errors.New("timed out waiting for connection: " + err.Error())
			}
		}
	}
}

func (r *ContainerDB) CreateRoles(usernames, passwords []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentDatabase string
	err = tx.Get(&currentDatabase, "SELECT current_database()")
	if err != nil {
		return err
	}

	var databaseToRevoke []string
	err = tx.Select(&databaseToRevoke,
		`SELECT datname FROM pg_database WHERE datistemplate = false AND datname != $1`,
		currentDatabase)

	queryRevokeConnect := `REVOKE CONNECT ON DATABASE "%s" FROM %s;`
	queryRevokeCreate := `REVOKE CREATE ON DATABASE "%s" FROM %s;`

	queryCreateRole := `CREATE ROLE %s WITH NOSUPERUSER NOCREATEDB NOCREATEROLE LOGIN PASSWORD '%s';`
	queryFunctions := `ALTER DEFAULT PRIVILEGES FOR ROLE %s IN SCHEMA public GRANT ALL PRIVILEGES ON FUNCTIONS TO PUBLIC;`

	queryAdminTables := `ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO admin;`
	queryAdminRedactorTables := `ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO redactor;`
	queryAdminReaderTables := `ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public GRANT SELECT ON TABLES TO reader;`
	queryAdminSchema := `GRANT ALL PRIVILEGES ON SCHEMA "public" TO %s;`

	queryRedactorAdminTables := `ALTER DEFAULT PRIVILEGES FOR ROLE redactor IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO admin;`
	queryRedactorTables := `ALTER DEFAULT PRIVILEGES FOR ROLE redactor IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO redactor;`
	queryRedactorReaderTables := `ALTER DEFAULT PRIVILEGES FOR ROLE redactor IN SCHEMA public GRANT SELECT ON TABLES TO reader;`
	queryRedactorSchema := `GRANT ALL PRIVILEGES ON SCHEMA "public" TO %s;`
	//queryRedactorSchema := `REVOKE ALL PRIVILEGES ON SCHEMA "public" FROM %s;`

	queryReaderAdminTables := `ALTER DEFAULT PRIVILEGES FOR ROLE reader IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO admin;`
	queryReaderRedactorTables := `ALTER DEFAULT PRIVILEGES FOR ROLE reader IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO redactor;`
	queryReaderTables := `ALTER DEFAULT PRIVILEGES FOR ROLE reader IN SCHEMA public GRANT SELECT ON TABLES TO reader;`
	queryReaderSchema := `GRANT ALL PRIVILEGES ON SCHEMA "public" TO %s;`

	for i, username := range usernames {
		if username == model.OwnerRole {
			continue
		}
		_, err = tx.Exec(fmt.Sprintf(queryCreateRole, username, passwords[i]))
		if err != nil {
			return err
		}
		_, err = tx.Exec(fmt.Sprintf(queryRevokeCreate, currentDatabase, username))
		if err != nil {
			return err
		}
		_, err = tx.Exec(fmt.Sprintf(queryFunctions, username))
		if err != nil {
			return err
		}
		for _, database := range databaseToRevoke {
			_, err = tx.Exec(fmt.Sprintf(queryRevokeConnect, database, username))
			if err != nil {
				return err
			}
		}
	}

	for _, username := range usernames {
		if username == model.OwnerRole {
			continue
		}
		var queryTables1, queryTables2, queryTables3, querySchema string

		switch username {
		case model.AdminRole:
			queryTables1 = queryAdminTables
			queryTables2 = queryAdminRedactorTables
			queryTables3 = queryAdminReaderTables
			querySchema = queryAdminSchema
		case model.RedactorRole:
			queryTables1 = queryRedactorAdminTables
			queryTables2 = queryRedactorTables
			queryTables3 = queryRedactorReaderTables
			querySchema = queryRedactorSchema
		case model.ReaderRole:
			queryTables1 = queryReaderAdminTables
			queryTables2 = queryReaderRedactorTables
			queryTables3 = queryReaderTables
			querySchema = queryReaderSchema
		}

		_, err = tx.Exec(queryTables1)
		if err != nil {
			return err
		}
		_, err = tx.Exec(queryTables2)
		if err != nil {
			return err
		}
		_, err = tx.Exec(queryTables3)
		if err != nil {
			return err
		}
		_, err = tx.Exec(fmt.Sprintf(querySchema, username))
		if err != nil {
			return err
		}

	}

	return tx.Commit()
}
