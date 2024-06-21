package myerr

import (
	"dbb-server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func New(c *gin.Context, statusCode int, message string) {
	logrus.Error(statusCode, " method: ", c.Request.Method, ", url: ", c.Request.URL, ", msg: ", message)

	c.AbortWithStatusJSON(statusCode, model.Response{
		Status:  statusCode,
		Message: message,
		Payload: nil,
	})
}

func NewErrorWithType(c *gin.Context, err error) {
	switch err.(type) {
	case BadRequest:
		New(c, http.StatusBadRequest, err.Error())
	case InternalError:
		New(c, http.StatusInternalServerError, err.Error())
	case Forbidden:
		New(c, http.StatusForbidden, err.Error())
	default:
		New(c, http.StatusTeapot, err.Error())
	}
	return
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
