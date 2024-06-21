package service

import (
	"dbb-server/internal/dockercli"
	"dbb-server/internal/model"
	"dbb-server/internal/repository"
)

type Auth interface {
	CreateUser(login, password string) (int, error)
	GetUserByCredentials(login, password string) (model.User, error)

	GenerateTokensAndSave(userId int) (string, string, error)
	GenerateTokensPair(userId int) (string, string, error)
	ParseAccessToken(accessToken string) (*model.AccessTokenClaimsExtension, error)
	ParseRefreshToken(refreshToken string) (*model.RefreshTokenClaimsExtension, error)
	RegenerateTokens(userId int, refreshToken, jti string) (string, string, error)

	DeleteSessionForUser(userId int) error
}

type Organization interface {
	GetOrganizationsForUser(userId int) ([]model.OrganizationForUser, error)
	CreteOrganization(userId int, name string) (int, error)
	GetUserRoleInOrganization(userId, organizationId int) (string, error)
	DeleteOrganization(organizationId int) (int, error)
	ChangeOrganizationName(organizationId int, name string) (int, error)

	InviteUserToOrganization(organizationId, userId int, role string) error
	DeleteUserFromOrganization(userId, organizationId int) (int, error)
	ChangeUserRoleInOrganization(userId, organizationId int, role string) (int, error)
}

type User interface {
	GetAllUsersInOrganization(organizationId, limit, page int, role string) (int, []model.UserInOrganization, error)
	GetAllUsers(limit, page int, search string) (int, []model.UserWithoutPassword, error)
}

type Datasource interface {
	CreateDataSource(organizationId int, dbName string) (int, error)
	GetDatasourcesInOrganization(organizationId, userId int) ([]model.DatasourceInOrganization, error)
	DeleteDatasource(datasourceId, userId int) (int, error)

	GetDatasourceData(datasourceId, userId int) (model.Datasource, model.DatasourceUser, error)
}

type Service struct {
	Auth
	Organization
	User
	Datasource
}

func NewService(repo *repository.Repository, cli *dockercli.DockerClient) *Service {
	return &Service{
		Auth:         NewAuthService(repo.Auth, cli),
		Organization: NewOrganizationService(repo.Organization, cli),
		User:         NewUserService(repo.User, cli),
		Datasource:   NewDataSourceService(repo.Datasource, cli),
	}
}
