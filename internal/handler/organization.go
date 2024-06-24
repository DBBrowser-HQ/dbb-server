package handler

import (
	"dbb-server/internal/model"
	"dbb-server/internal/myerr"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllUsersOrganizations(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	orgs, err := h.services.Organization.GetOrganizationsForUser(userData.UserId)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: orgs,
	})
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err = c.ShouldBindJSON(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Organization.CreteOrganization(userData.UserId, input.Name)
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

func (h *Handler) DeleteOrganization(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	organizationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	userRole, err := h.services.Organization.GetUserRoleInOrganization(userData.UserId, organizationId)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userRole != model.AdminRole {
		myerr.New(c, http.StatusForbidden, "You are not an admin in this organization")
		return
	}

	id, err := h.services.Organization.DeleteOrganization(organizationId)
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

func (h *Handler) ChangeOrganizationName(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	organizationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err = c.ShouldBind(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	userRole, err := h.services.Organization.GetUserRoleInOrganization(userData.UserId, organizationId)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userRole != model.AdminRole {
		myerr.New(c, http.StatusForbidden, "You are not an admin in this organization")
		return
	}

	id, err := h.services.Organization.ChangeOrganizationName(organizationId, input.Name)
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

func (h *Handler) InviteUserToOrganization(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	newUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	if newUserId == userData.UserId {
		myerr.New(c, http.StatusBadRequest, "You can't invite yourself")
		return
	}

	var input struct {
		OrganizationId int    `json:"organizationId" binding:"required"`
		Role           string `json:"role" binding:"required"`
	}

	if err = c.ShouldBind(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Role != model.AdminRole && input.Role != model.RedactorRole && input.Role != model.ReaderRole {
		myerr.New(c, http.StatusBadRequest, "Unknown role")
		return
	}

	userRole, err := h.services.Organization.GetUserRoleInOrganization(userData.UserId, input.OrganizationId)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userRole != model.AdminRole {
		myerr.New(c, http.StatusForbidden, "You are not an admin in this organization")
		return
	}

	if err = h.services.Organization.InviteUserToOrganization(input.OrganizationId, newUserId, input.Role); err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) DeleteUserFromOrganization(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	newUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	if newUserId == userData.UserId {
		myerr.New(c, http.StatusBadRequest, "Don't do it please")
		return
	}

	var input struct {
		OrganizationId int `json:"organizationId" binding:"required"`
	}

	if err = c.ShouldBind(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	userRole, err := h.services.Organization.GetUserRoleInOrganization(userData.UserId, input.OrganizationId)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userRole != model.AdminRole {
		myerr.New(c, http.StatusForbidden, "You are not an admin in this organization")
		return
	}

	id, err := h.services.Organization.DeleteUserFromOrganization(newUserId, input.OrganizationId)
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

func (h *Handler) ChangeUserRoleInOrganization(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	newUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	if newUserId == userData.UserId {
		myerr.New(c, http.StatusBadRequest, "Don't do it please")
		return
	}

	var input struct {
		OrganizationId int    `json:"organizationId" binding:"required"`
		Role           string `json:"role" binding:"required"`
	}

	if err = c.ShouldBind(&input); err != nil {
		myerr.New(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Role != model.AdminRole && input.Role != model.RedactorRole && input.Role != model.ReaderRole {
		myerr.New(c, http.StatusBadRequest, "Unknown role")
		return
	}

	userRole, err := h.services.Organization.GetUserRoleInOrganization(userData.UserId, input.OrganizationId)
	if err != nil {
		myerr.New(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userRole != model.AdminRole {
		myerr.New(c, http.StatusForbidden, "You are not an admin in this organization")
		return
	}

	id, err := h.services.Organization.ChangeUserRoleInOrganization(newUserId, input.OrganizationId, input.Role)
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
