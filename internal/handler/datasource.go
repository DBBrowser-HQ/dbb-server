package handler

import (
	"dbb-server/internal/model"
	"dbb-server/internal/myerr"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateDatasource(c *gin.Context) {
	var request struct {
		Name           string `json:"name" binding:"required"`
		OrganizationId int    `json:"organizationId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Datasource.CreateDataSource(request.OrganizationId, request.Name)
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

func (h *Handler) DeleteDatasource(c *gin.Context) {

}
