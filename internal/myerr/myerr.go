package myerr

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

type BadRequest struct {
	Err string `json:"error"`
}

func NewBadRequest(err string) BadRequest {
	return BadRequest{Err: err}
}

func (b BadRequest) Error() string {
	return b.Err
}

type InternalError struct {
	Err string `json:"error"`
}

func NewInternalError(err string) InternalError {
	return InternalError{Err: err}
}

func (b InternalError) Error() string {
	return b.Err
}

type Forbidden struct {
	Err string `json:"error"`
}

func NewForbidden(err string) Forbidden {
	return Forbidden{Err: err}
}

func (b Forbidden) Error() string {
	return b.Err
}
