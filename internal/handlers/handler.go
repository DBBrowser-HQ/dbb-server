package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/refresh", h.RefreshToken)
	}

	api := router.Group("/api", h.UserIdentify)
	{
		orgs := api.Group("/orgs")
		{
			orgs.GET("", h.GetAllUsersOrganizations)
			orgs.POST("", h.CreateOrganization)
			orgs.DELETE("/:id", h.DeleteOrganization)
			orgs.PATCH(":/id", h.ChangeOrganizationName)

			orgsUsers := orgs.Group("")
			{
				orgsUsers.POST("/invite/:id", h.InviteUserToOrganization)
				orgsUsers.POST("/kick/:id", h.DeleteUserFromOrganization)
				orgsUsers.POST("/change-role/:id", h.ChangeUserRoleInOrganization)
			}
		}
		users := api.Group("/users")
		{
			users.GET("", h.GetAllUsers)
		}
		datasources := api.Group("/datasources")
		{
			datasources.POST("/:id", h.CreateDatasource)
			datasources.DELETE("/:id", h.DeleteDatasource)
		}
	}

	return router
}
