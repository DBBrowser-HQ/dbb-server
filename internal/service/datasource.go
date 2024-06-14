package service

import (
	"context"
	"dbb-server/internal/dockercli"
	"dbb-server/internal/model"
	"dbb-server/internal/repository"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jmrobles/h2go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"os"
	"time"
)

type DataSourceService struct {
	repo repository.Datasource
	cli  *dockercli.DockerClient
}

func NewDataSourceService(repo repository.Datasource, cli *dockercli.DockerClient) *DataSourceService {
	return &DataSourceService{repo: repo, cli: cli}
}

func generateHash() string {
	id := uuid.New().String()

	hash := sha3.New256()
	hash.Write([]byte(id))

	salt := os.Getenv("HASH_SALT")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *DataSourceService) CreateDataSource(organizationId int, dbName string) (int, error) {
	exist, err := s.repo.CheckDatasourceExistence(dbName, organizationId)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, errors.New(fmt.Sprintf("datasource with name: %s already exists in this organization", dbName))
	}

	dbHost := os.Getenv("DB_HOST") + "-" + uuid.New().String()
	dbPort, err := s.repo.GetUnusedPort()
	if err != nil {
		return 0, err
	}
	dbUsername := model.OwnerRole
	dbPassword := generateHash()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	actualDbPort, err := s.cli.CreateDockerContainer(ctx, dbHost, dbPort, dbUsername, dbPassword, dbName)
	if err != nil {
		return 0, err
	}

	usernames := make([]string, 0)
	passwords := make([]string, 0)
	usernames = append(usernames, dbUsername, model.AdminRole, model.RedactorRole, model.ReaderRole)
	passwords = append(passwords, dbPassword, generateHash(), generateHash(), generateHash())

	db, err := repository.NewContainerDB(dbHost, dbUsername, dbPassword, dbName, actualDbPort)
	if err != nil {
		return 0, err
	}

	if err = db.CreateRoles(usernames, passwords); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateDatasource(dbHost, dbName, actualDbPort, organizationId)
	if err != nil {
		logrus.Info("here 3")
		return 0, err
	}

	if err = s.repo.AddRolesData(usernames, passwords, id); err != nil {
		return 0, err
	}

	return id, nil
}
