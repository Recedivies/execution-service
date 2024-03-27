package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/config"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/exception"
	"gitlab.cs.ui.ac.id/ahmadhi.prananta/execution_service/util"
	"gorm.io/gorm"
)

// Server serves HTTP requests.
type Server struct {
	config config.Config
	DB     *gorm.DB
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(configuration config.Config, db *gorm.DB) (*Server, error) {
	server := &Server{
		config: configuration,
		DB:     db,
	}

	if server.config.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()
	router.Use(util.CORSMiddleware())
	router.Use(util.HttpLogger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) start(address string) error {
	return server.router.Run(address)
}

func RunGinServer(configuration config.Config, db *gorm.DB) {
	server, err := NewServer(configuration, db)
	exception.FatalIfNeeded(err, "cannot create server")

	err = server.start(configuration.HTTPServerAddress)
	exception.FatalIfNeeded(err, "cannot start server")
}
