package service

import (
	"dbb-server/internal/dockercli"
	"dbb-server/internal/model"
	"dbb-server/internal/repository"
)

type UserService struct {
	repo repository.User
	cli  *dockercli.DockerClient
}

func NewUserService(repo repository.User, cli *dockercli.DockerClient) *UserService {
	return &UserService{repo: repo, cli: cli}
}

func (s *UserService) GetAllUsersInOrganization(organizationId, limit, page int, role string) (int, []model.UserInOrganization, error) {
	offset := limit * page

	return s.repo.GetAllUsersInOrganization(organizationId, limit, offset, role)
}

func (s *UserService) GetAllUsers(limit, page int, search string) (int, []model.UserWithoutPassword, error) {
	offset := limit * page

	return s.repo.GetAllUsers(limit, offset, search)
}
