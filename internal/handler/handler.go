package handler

import (
	"dbb-server/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Any("/connect/:dbName", h.ServeConnection)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/refresh", h.UserIdentify, h.RefreshTokens)
		auth.POST("/logout", h.UserIdentify, h.Logout)
	}

	api := router.Group("/api", h.UserIdentify)
	{
		organizations := api.Group("/organizations")
		{
			organizations.GET("", h.GetAllUsersOrganizations)
			organizations.POST("", h.CreateOrganization)
			organizations.DELETE("/:id", h.DeleteOrganization)
			organizations.PATCH("/:id", h.ChangeOrganizationName)

			organizationsUsers := organizations.Group("")
			{
				organizationsUsers.POST("/invite/:id", h.InviteUserToOrganization)
				organizationsUsers.POST("/kick/:id", h.DeleteUserFromOrganization)
				organizationsUsers.POST("/change-role/:id", h.ChangeUserRoleInOrganization)
			}
		}
		users := api.Group("/users")
		{
			users.GET("", h.GetAllUsers)
			users.GET("/:id", h.GetAllUsersInOrganization)
		}
		datasources := api.Group("/datasources")
		{
			datasources.POST("/:id", h.CreateDatasource)
			datasources.DELETE("/:id", h.DeleteDatasource)
		}
	}

	return router
}
