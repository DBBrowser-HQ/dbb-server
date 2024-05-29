package handler

import (
	"dbb-server/internal/myerrors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
)

// ServeConnection DON'T WORK
func (h *Handler) ServeConnection(c *gin.Context) {
	dbHost := os.Getenv("H2_DB_HOST")
	dbPort := os.Getenv("H2_DB_PORT")
	dbAddress := dbHost + ":" + dbPort

	dbConn, err := net.Dial("tcp", dbAddress)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer dbConn.Close()

	hijacker, ok := c.Writer.(http.Hijacker)
	if !ok {
		myerrors.New(c, http.StatusInternalServerError, "Hijacking error")
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer clientConn.Close()

	go func() {
		logrus.Info("Start messaging")
		if _, err = io.Copy(dbConn, clientConn); err != nil {
			logrus.Println("Error forwarding data from client to db:", err)
		}
	}()

	if _, err = io.Copy(clientConn, dbConn); err != nil {
		logrus.Println("Error forwarding data from db to client:", err)
	}
}
