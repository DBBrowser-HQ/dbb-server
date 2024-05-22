package handlers

import (
	"crypto/sha1"
	"database/sql"
	"dbb-server/internal/db_postgres"
	"dbb-server/internal/model"
	"dbb-server/internal/myerrors"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

const (
	accessTTL  = 3 * time.Hour
	refreshTTL = 60 * 24 * time.Hour
)

type accessTokenClaims struct {
	jwt.StandardClaims
	model.TokenClaimsExtension
}

type refreshTokenClaims struct {
	jwt.StandardClaims
	model.TokenClaimsExtension
}

func generatePasswordHash(login, password string) string {
	hash := sha1.New()
	hash.Write([]byte(login + password))

	salt := os.Getenv("PASSWORD_SALT")
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (h *Handler) SignUp(c *gin.Context) {
	var input model.SignUser
	if err := c.ShouldBindJSON(&input); err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash := generatePasswordHash(input.Login, input.Password)

	query := fmt.Sprintf(`INSERT INTO %s (login, password_hash) VALUES ($1, $2) RETURNING id`, db_postgres.UsersTable)
	var id int
	err := h.db.Get(&id, query, input.Login, passwordHash)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"id": id,
		},
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	var input model.SignUser
	if err := c.ShouldBindJSON(&input); err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	query := fmt.Sprintf(`SELECT id, login, password_hash FROM %s
                                WHERE login=$1 AND password_hash=$2`, db_postgres.UsersTable)

	var user model.User
	err := h.db.Get(&user, query, input.Login, generatePasswordHash(input.Login, input.Password))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			myerrors.New(c, http.StatusUnauthorized, "wrong credentials")
			return
		}
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := generateTokensPair(user.Id)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refreshToken" binging:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := ParseRefreshToken(input.RefreshToken)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := generateTokensPair(claims.UserId)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

func generateTokensPair(userId int) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &accessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenClaimsExtension: model.TokenClaimsExtension{
			UserId: userId,
		},
	})

	uuidVar := uuid.New().String()
	if uuidVar == "" {
		return "", "", errors.New("can't generate uuid for jti")
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuidVar,
		},
		TokenClaimsExtension: model.TokenClaimsExtension{
			UserId: userId,
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

func ParseAccessToken(accessToken string) (*model.TokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(accessToken, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("ACCESS_SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*accessTokenClaims); !ok {
		return nil, errors.New("invalid access token claims type")
	} else {
		return &claims.TokenClaimsExtension, nil
	}
}

func ParseRefreshToken(refreshToken string) (*model.TokenClaimsExtension, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("REFRESH_SIGNING_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*refreshTokenClaims); !ok {
		return nil, errors.New("invalid refresh token claims type")
	} else {
		return &claims.TokenClaimsExtension, nil
	}
}
