package handlers

import (
	"dbb-server/internal/myerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authHeader = "Authorization"
)

func (h *Handler) UserIdentify(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		myerrors.New(c, http.StatusUnauthorized, "Empty Authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		myerrors.New(c, http.StatusUnauthorized, "Invalid Authorization header")
		return
	}

	if len(headerParts[1]) == 0 {
		myerrors.New(c, http.StatusUnauthorized, "Empty token")
		return
	}

	userData, err := ParseAccessToken(headerParts[1])
	if err != nil {
		myerrors.New(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("userId", userData.UserId)
	c.Next()
}
