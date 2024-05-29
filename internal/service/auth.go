package service

import (
	"dbb-server/internal/dockercli"
	"dbb-server/internal/model"
	"dbb-server/internal/repository"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	"os"
	"time"
)

type AuthService struct {
	repo repository.Auth
	cli  *dockercli.DockerClient
}

func NewAuthService(repo repository.Auth, cli *dockercli.DockerClient) *AuthService {
	return &AuthService{repo: repo, cli: cli}
}

const (
	maxSessionsNumber = 5
	accessTokenTTL    = 1 * time.Hour
	refreshTokenTTL   = 7 * 24 * time.Hour
)

func generatePasswordHash(login, password string) string {
	hash := sha3.New256()
	hash.Write([]byte(login + password))

	salt := os.Getenv("PASSWORD_SALT")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) CreateUser(login, password string) (int, error) {
	passwordHash := generatePasswordHash(login, password)
	return s.repo.CreateUser(login, passwordHash)
}

func (s *AuthService) GetUserByCredentials(login, password string) (model.User, error) {
	passwordHash := generatePasswordHash(login, password)
	return s.repo.GetUserByCredentials(login, passwordHash)
}

func (s *AuthService) GenerateTokensPair(userId int) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		AccessTokenClaimsExtension: model.AccessTokenClaimsExtension{
			UserId: userId,
		},
	})

	jti := uuid.New().String()
	if jti == "" {
		return "", "", errors.New("can't generate uuid for jti")
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		RefreshTokenClaimsExtension: model.RefreshTokenClaimsExtension{
			UserId: userId,
			JTI:    jti,
		},
	})

	accessSigningKey := os.Getenv("ACCESS_SIGNING_KEY")
	signedAccessToken, err := accessToken.SignedString([]byte(accessSigningKey))
	if err != nil {
		return "", "", err
	}

	refreshSigningKey := os.Getenv("REFRESH_SIGNING_KEY")
	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshSigningKey))
	if err != nil {
		return "", "", err
	}
	return signedAccessToken, signedRefreshToken, nil
}

func (s *AuthService) GenerateTokensAndSave(userId int) (string, string, error) {
	accessToken, refreshToken, err := s.GenerateTokensPair(userId)
	if err != nil {
		return "", "", err
	}

	claims, err := s.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	err = s.repo.SaveSessionData(refreshToken, claims.JTI, userId)

	return accessToken, refreshToken, err
}

func (s *AuthService) ParseAccessToken(accessToken string) (*model.AccessTokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("ACCESS_SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.AccessTokenClaims); !ok {
		return nil, errors.New("invalid access token claims type")
	} else {
		return &claims.AccessTokenClaimsExtension, nil
	}
}

func (s *AuthService) ParseRefreshToken(refreshToken string) (*model.RefreshTokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &model.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("REFRESH_SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.RefreshTokenClaims); !ok {
		return nil, errors.New("invalid refresh token claims type")
	} else {
		return &claims.RefreshTokenClaimsExtension, nil
	}
}

func (s *AuthService) RegenerateTokens(userId int, refreshToken, jti string) (string, string, error) {
	sessionData, err := s.repo.GetSessionData(userId)
	if err != nil {
		return "", "", err
	}

	if sessionData.JTI != jti || sessionData.RefreshToken != refreshToken {
		err = s.DeleteSessionForUser(userId)
		if err != nil {
			return "", "", errors.New("suspicious activity detected, but: " + err.Error())
		}
		return "", "", errors.New("suspicious activity")
	}

	accessToken, refreshToken, err := s.GenerateTokensPair(userId)
	if err != nil {
		return "", "", err
	}

	claims, err := s.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	err = s.repo.UpdateSessionData(refreshToken, claims.JTI, userId)
	return accessToken, refreshToken, err
}

func (s *AuthService) DeleteSessionForUser(userId int) error {
	return s.repo.DeleteSessionForUser(userId)
}
