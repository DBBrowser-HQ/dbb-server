package repository

import (
	"dbb-server/internal/model"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(login, passwordHash string) (int, error)
	GetUserByCredentials(login, password string) (model.User, error)

	SaveSessionData(refreshToken, jti string, userId int) error
	GetSessionData(userId int) (model.RefreshSession, error)

	DeleteSessionForUser(userId int) error
	UpdateSessionData(refreshToken, jti string, userId int) error
}

type Organization interface {
	GetOrganizationsForUser(userId int) ([]model.OrganizationForUser, error)
	CreateOrganization(name string, userId int) (int, error)
	AddUserWithRoleToOrganization(userId, organizationId int, role string) error
	GetUserRoleInOrganization(userId, organizationId int) (string, error)
	DeleteOrganization(organizationId int) (int, error)
	ChangeOrganizationName(organizationId int, name string) (int, error)
	DeleteUserFromOrganization(userId, organizationId int) (int, error)
	ChangeUserRoleInOrganization(userId, organizationId int, role string) (int, error)
}

type User interface {
	GetAllUsersInOrganization(organizationId, limit, offset int, role string) (int, []model.UserInOrganization, error)
	GetAllUsers(limit, offset int, search string) (int, []model.UserWithoutPassword, error)
}

type Datasource interface {
	GetUser(userId int) (string, string, error)
	GetUnusedPort() (int, error)
	CheckDatasourceExistence(dbName string, organizationId int) (bool, error)
	CreateDatasource(dbHost, dbName string, dbPort, organizationId int) (int, error)
	AddRolesData(usernames, passwords []string, datasourceId int) error
	GetHostNames() ([]string, error)

	CheckUserInOrganization(organizationId, userId int) (bool, error)
	GetDatasourcesInOrganization(organizationId int) ([]model.DatasourceInOrganization, error)
	GetDatasourceData(datasourceId int, role string) (model.Datasource, model.DatasourceUser, error)

	GetUserRoleInDatasource(datasourceId, userId int) (string, error)
	DeleteDatasource(datasourceId int) (int, string, error)
}

type Repository struct {
	Auth
	Organization
	User
	Datasource
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:         NewAuthPostgres(db),
		Organization: NewOrganizationPostgres(db),
		User:         NewUserPostgres(db),
		Datasource:   NewDataSourcePostgres(db),
	}
}
