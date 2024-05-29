package repository

import (
	"dbb-server/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetAllUsersInOrganization(organizationId, limit, offset int, role string) (int, []model.UserInOrganization, error) {
	queryCount := fmt.Sprintf(`SELECT COUNT(*)
									FROM %s o
									JOIN %s u
									ON o.user_id = u.id
									WHERE o.organization_id = $1`, UsersOrganizationsTable, UsersTable)

	queryUsers := fmt.Sprintf(`SELECT u.id, u.login, o.role
									FROM %s o
									JOIN %s u
									ON o.user_id = u.id
									WHERE o.organization_id = $1`, UsersOrganizationsTable, UsersTable)

	if role != "" {
		queryUsers += ` AND o.role = $2 LIMIT $3 OFFSET $4`
	} else {
		queryUsers += ` LIMIT $2 OFFSET $3`
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	var count int
	err = tx.Get(&count, queryCount, organizationId)
	if err != nil {
		return 0, nil, err
	}

	var users []model.UserInOrganization
	if role != "" {
		err = tx.Select(&users, queryUsers, organizationId, role, limit, offset)
	} else {
		err = tx.Select(&users, queryUsers, organizationId, limit, offset)
	}
	if err != nil {
		return 0, nil, err
	}
	err = tx.Commit()
	return count, users, err
}

func (r *UserPostgres) GetAllUsers(limit, offset int, search string) (int, []model.UserWithoutPassword, error) {
	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, UsersTable)

	queryUsers := fmt.Sprintf(`SELECT id, login FROM %s`, UsersTable)

	if search != "" {
		queryUsers += ` WHERE login LIKE '%' || $1 || '%' LIMIT $2 OFFSET $3`
	} else {
		queryUsers += ` LIMIT $1 OFFSET $2`
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	var count int
	err = tx.Get(&count, queryCount)
	if err != nil {
		return 0, nil, err
	}

	var users []model.UserWithoutPassword
	if search != "" {
		err = tx.Select(&users, queryUsers, search, limit, offset)
	} else {
		err = tx.Select(&users, queryUsers, limit, offset)
	}
	if err != nil {
		return 0, nil, err
	}
	err = tx.Commit()
	return count, users, err
}
