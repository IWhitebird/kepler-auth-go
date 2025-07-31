package api

import (
	"kepler-auth-go/internal/config"
	"kepler-auth-go/internal/handlers"
	"kepler-auth-go/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	cfg                 *config.Config
	authHandler         *handlers.AuthHandler
	userHandler         *handlers.UserHandler
	emailHandler        *handlers.EmailHandler
	groupHandler        *handlers.GroupHandler
	permissionHandler   *handlers.PermissionHandler
	organizationHandler *handlers.OrganizationHandler
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:                 cfg,
		authHandler:         handlers.NewAuthHandler(cfg),
		userHandler:         handlers.NewUserHandler(),
		emailHandler:        handlers.NewEmailHandler(cfg),
		groupHandler:        handlers.NewGroupHandler(),
		permissionHandler:   handlers.NewPermissionHandler(),
		organizationHandler: handlers.NewOrganizationHandler(),
	}
}

func (s *Server) SetupRouter() *gin.Engine {
	gin.SetMode(s.cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s.setupRoutes(r)

	return r
}
