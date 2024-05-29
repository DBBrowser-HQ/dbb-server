package handler

import (
	"dbb-server/internal/model"
	"dbb-server/internal/myerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllUsers(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	search := c.Query("search")

	count, users, err := h.services.User.GetAllUsers(limitInt, pageInt, search)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": count,
			"rows":  users,
		},
	})

}

func (h *Handler) GetAllUsersInOrganization(c *gin.Context) {
	organizationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		myerrors.New(c, http.StatusBadRequest, err.Error())
		return
	}

	role := c.Query("role")
	if role != "" {
		if role != model.AdminRole && role != model.RedactorRole && role != model.ReaderRole {
			myerrors.New(c, http.StatusBadRequest, "unknown role")
			return
		}
	}

	count, users, err := h.services.User.GetAllUsersInOrganization(organizationId, limitInt, pageInt, role)
	if err != nil {
		myerrors.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": count,
			"rows":  users,
		},
	})
}
