package service

import (
	"dbb-server/internal/dockercli"
	"dbb-server/internal/model"
	"dbb-server/internal/repository"
)

type OrganizationService struct {
	repo repository.Organization
	cli  *dockercli.DockerClient
}

func NewOrganizationService(repo repository.Organization, cli *dockercli.DockerClient) *OrganizationService {
	return &OrganizationService{repo: repo, cli: cli}
}

func (s *OrganizationService) GetOrganizationsForUser(userId int) ([]model.OrganizationForUser, error) {
	return s.repo.GetOrganizationsForUser(userId)
}

func (s *OrganizationService) CreteOrganization(userId int, name string) (int, error) {
	id, err := s.repo.CreateOrganization(name, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *OrganizationService) GetUserRoleInOrganization(userId, organizationId int) (string, error) {
	return s.repo.GetUserRoleInOrganization(userId, organizationId)
}

func (s *OrganizationService) DeleteOrganization(organizationId int) (int, error) {
	return s.repo.DeleteOrganization(organizationId)
}

func (s *OrganizationService) ChangeOrganizationName(organizationId int, name string) (int, error) {
	return s.repo.ChangeOrganizationName(organizationId, name)
}

func (s *OrganizationService) InviteUserToOrganization(organizationId, userId int, role string) error {
	return s.repo.AddUserWithRoleToOrganization(userId, organizationId, role)
}

func (s *OrganizationService) DeleteUserFromOrganization(userId, organizationId int) (int, error) {
	return s.repo.DeleteUserFromOrganization(userId, organizationId)
}

func (s *OrganizationService) ChangeUserRoleInOrganization(userId, organizationId int, role string) (int, error) {
	return s.repo.ChangeUserRoleInOrganization(userId, organizationId, role)
}
