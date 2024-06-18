package handler

import (
	"dbb-server/internal/model"
	"dbb-server/internal/myerr"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) CreateDatasource(c *gin.Context) {
	orgId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Datasource.CreateDataSource(orgId, request.Name)
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

func (h *Handler) GetDatasourcesInOrganization(c *gin.Context) {
	orgId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	datasources, err := h.services.Datasource.GetDatasourcesInOrganization(orgId, userData.UserId)
	if err != nil {
		switch err.(type) {
		case myerr.BadRequest:
			myerr.New(c, http.StatusBadRequest, err.Error())
		case myerr.InternalError:
			myerr.New(c, http.StatusInternalServerError, err.Error())
		default:
			myerr.New(c, http.StatusTeapot, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: datasources,
	})
}

func (h *Handler) DeleteDatasource(c *gin.Context) {
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

	id, err := h.services.Datasource.DeleteDatasource(datasourceId, userData.UserId)
	if err != nil {
		switch err.(type) {
		case myerr.BadRequest:
			myerr.New(c, http.StatusBadRequest, err.Error())
		case myerr.InternalError:
			myerr.New(c, http.StatusInternalServerError, err.Error())
		case myerr.Forbidden:
			myerr.New(c, http.StatusForbidden, err.Error())
		default:
			myerr.New(c, http.StatusTeapot, err.Error())
		}
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
