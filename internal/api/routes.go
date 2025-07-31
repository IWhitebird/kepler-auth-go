package api

import (
	"kepler-auth-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		s.setupAuthRoutes(api)
		s.setupUserRoutes(api)
		s.setupEmailRoutes(api)
		s.setupGroupRoutes(api)
		s.setupPermissionRoutes(api)
		s.setupOrganizationRoutes(api)
	}
}

func (s *Server) setupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", s.authHandler.Register)
		auth.POST("/login", s.authHandler.Login)

		authenticated := auth.Group("/")
		authenticated.Use(middleware.AuthRequired(s.cfg))
		{
			authenticated.GET("/me", s.authHandler.GetMe)
			authenticated.PATCH("/me", s.authHandler.UpdateMe)
			authenticated.POST("/change-password", s.authHandler.ChangePassword)
		}
	}
}

func (s *Server) setupUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	users.Use(middleware.AuthRequired(s.cfg))
	{
		users.GET("", s.userHandler.GetUsers)
		users.GET("/:id", s.userHandler.GetUser)

		adminRequired := users.Group("/")
		adminRequired.Use(middleware.AdminRequired())
		{
			adminRequired.PATCH("/:id", s.userHandler.UpdateUser)
			adminRequired.DELETE("/:id", s.userHandler.DeleteUser)
		}
	}
}

func (s *Server) setupEmailRoutes(api *gin.RouterGroup) {
	email := api.Group("/email")
	email.Use(middleware.AuthRequired(s.cfg))
	{
		email.POST("/send", s.emailHandler.SendEmail)
	}
}

func (s *Server) setupGroupRoutes(api *gin.RouterGroup) {
	groups := api.Group("/groups")
	groups.Use(middleware.AuthRequired(s.cfg))
	{
		groups.GET("", s.groupHandler.GetGroups)
		groups.GET("/:id", s.groupHandler.GetGroup)

		adminRequired := groups.Group("/")
		adminRequired.Use(middleware.AdminRequired())
		{
			adminRequired.POST("", s.groupHandler.CreateGroup)
			adminRequired.PATCH("/:id", s.groupHandler.UpdateGroup)
			adminRequired.DELETE("/:id", s.groupHandler.DeleteGroup)
		}
	}
}

func (s *Server) setupPermissionRoutes(api *gin.RouterGroup) {
	permissions := api.Group("/permissions")
	permissions.Use(middleware.AuthRequired(s.cfg))
	{
		permissions.GET("", s.permissionHandler.GetPermissions)
	}

	authGroups := api.Group("/auth-groups")
	authGroups.Use(middleware.AuthRequired(s.cfg))
	{
		authGroups.GET("", s.permissionHandler.GetAuthGroups)
		authGroups.GET("/:id", s.permissionHandler.GetAuthGroup)

		adminRequired := authGroups.Group("/")
		adminRequired.Use(middleware.AdminRequired())
		{
			adminRequired.POST("", s.permissionHandler.CreateAuthGroup)
			adminRequired.PATCH("/:id", s.permissionHandler.UpdateAuthGroup)
			adminRequired.DELETE("/:id", s.permissionHandler.DeleteAuthGroup)
		}
	}
}

func (s *Server) setupOrganizationRoutes(api *gin.RouterGroup) {
	organizations := api.Group("/organizations")
	organizations.Use(middleware.AuthRequired(s.cfg))
	organizations.Use(middleware.AdminRequired()) // Only admins can manage organizations
	{
		organizations.GET("", s.organizationHandler.GetOrganizations)
		organizations.GET("/:id", s.organizationHandler.GetOrganization)
		organizations.POST("", s.organizationHandler.CreateOrganization)
		organizations.PATCH("/:id", s.organizationHandler.UpdateOrganization)
		organizations.DELETE("/:id", s.organizationHandler.DeleteOrganization)
	}
}
