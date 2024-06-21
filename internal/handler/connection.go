package handler

import (
	"dbb-server/internal/myerr"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strconv"
)

func findDbbProxyIp() string {
	ips, _ := net.LookupIP("dbb-proxy")
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}

var (
	dbbProxyIp = findDbbProxyIp()
)

func (h *Handler) ServeConnection(c *gin.Context) {
	if c.ClientIP() != dbbProxyIp {
		if dbbProxyIp == "" {
			myerr.New(c, http.StatusInternalServerError, "Can't find dbb-proxy ip")
			return
		}
		myerr.New(c, http.StatusForbidden, "only for proxy")
		return
	}

	datasourceId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	datasource, user, err := h.services.Datasource.GetDatasourceData(datasourceId, userData.UserId)
	if err != nil {
		myerr.NewErrorWithType(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host":     datasource.Host,
		"port":     datasource.Port,
		"user":     user.Username,
		"password": user.Password,
		"name":     datasource.Name,
	})
}
