package myerrors

import (
	"dbb-server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func New(c *gin.Context, statusCode int, message string) {
	logrus.Error(statusCode, " method: ", c.Request.Method, ", url: ", c.Request.URL, ", msg: ", message)

	c.AbortWithStatusJSON(statusCode, model.Response{
		Status:  statusCode,
		Message: message,
		Payload: nil,
	})
}
