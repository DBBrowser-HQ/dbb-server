package repository

import (
	"dbb-server/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type OrganizationPostgres struct {
	db *sqlx.DB
}

func NewOrganizationPostgres(db *sqlx.DB) *OrganizationPostgres {
	return &OrganizationPostgres{db: db}
}

func (r *OrganizationPostgres) GetOrganizationsForUser(userId int) ([]model.OrganizationForUser, error) {
	query := fmt.Sprintf(`SELECT o.id, o.name, u.role
								FROM %s o
								JOIN %s u
								ON o.id = u.organization_id
								WHERE u.user_id=$1`, OrganizationsTable, UsersOrganizationsTable)

	var organizations []model.OrganizationForUser
	err := r.db.Select(&organizations, query, userId)
	return organizations, err
}

func (r *OrganizationPostgres) CreateOrganization(name string, userId int) (int, error) {
	queryCreateOrg := fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1) RETURNING id`, OrganizationsTable)
	queryAddUser := fmt.Sprintf(`INSERT INTO %s (user_id, role, organization_id) VALUES ($1, $2, $3)`, UsersOrganizationsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	err = tx.Get(&id, queryCreateOrg, name)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(queryAddUser, userId, model.AdminRole, id)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *OrganizationPostgres) AddUserWithRoleToOrganization(userId, organizationId int, role string) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, role, organization_id) VALUES ($1, $2, $3)`, UsersOrganizationsTable)

	_, err := r.db.Exec(query, userId, role, organizationId)
	return err
}

func (r *OrganizationPostgres) GetUserRoleInOrganization(userId, organizationId int) (string, error) {
	query := fmt.Sprintf(`SELECT role FROM %s WHERE user_id=$1 AND organization_id=$2`, UsersOrganizationsTable)

	var role string
	err := r.db.Get(&role, query, userId, organizationId)
	return role, err
}

func (r *OrganizationPostgres) DeleteOrganization(organizationId int) (int, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1 RETURNING id`, OrganizationsTable)

	var id int
	err := r.db.Get(&id, query, organizationId)
	return id, err
}

func (r *OrganizationPostgres) ChangeOrganizationName(organizationId int, name string) (int, error) {
	query := fmt.Sprintf(`UPDATE %s SET name=$1 WHERE id=$2 RETURNING id`, OrganizationsTable)

	var id int
	err := r.db.Get(&id, query, name, organizationId)

	return id, err
}

func (r *OrganizationPostgres) DeleteUserFromOrganization(userId, organizationId int) (int, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1 AND organization_id=$2 RETURNING user_id`, UsersOrganizationsTable)

	var id int
	err := r.db.Get(&id, query, userId, organizationId)
	return id, err
}

func (r *OrganizationPostgres) ChangeUserRoleInOrganization(userId, organizationId int, role string) (int, error) {
	query := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE user_id=$2 AND organization_id=$3 RETURNING user_id`, UsersOrganizationsTable)

	var id int
	err := r.db.Get(&id, query, role, userId, organizationId)
	return id, err
}
