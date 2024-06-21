package handler

import (
	"database/sql"
	"dbb-server/internal/model"
	"dbb-server/internal/myerr"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) SignIn(c *gin.Context) {
	var input model.SignUser
	if err := c.ShouldBindJSON(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.Auth.GetUserByCredentials(input.Login, input.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			myerr.New(c, http.StatusUnauthorized, "wrong credentials: user wasn't found")
			return
		}
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Auth.GenerateTokensAndSave(user.Id)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
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

func (h *Handler) SignUp(c *gin.Context) {
	var input model.SignUser
	if err := c.ShouldBindJSON(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Auth.CreateUser(input.Login, input.Password)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
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

func (h *Handler) RefreshTokens(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := h.services.Auth.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		myerr.NewErrorWithType(c, err)
		return
	}

	accessToken, refreshToken, err := h.services.Auth.RegenerateTokens(claims.UserId, input.RefreshToken, claims.JTI)
	if err != nil {
		if strings.Contains(err.Error(), "suspicious activity") {
			myerr.New(c, http.StatusUnauthorized, err.Error())
			return
		}
		myerr.New(c, http.StatusInternalServerError, err.Error())
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

func (h *Handler) Logout(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Auth.DeleteSessionForUser(userData.UserId); err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}
