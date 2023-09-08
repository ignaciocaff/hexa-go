package server

import (
	"context"
	"fmt"
	"support/internal/env"
	"support/internal/handlers"

	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

type Starter interface {
	Start()
}

type AppServer struct {
	*gin.Engine
	ctx    context.Context
	db     *sqlx.DB
	config env.EnvApp
}

func (s *AppServer) configure() {
	s.Engine.Use(gin.Recovery())
	s.Use(handlers.StatsMiddleware())
	s.Use(handlers.RestartOnErrorMiddleware())

	s.Engine.SetTrustedProxies([]string{"*"})
	// use CORS middleware
}

func NewGinServer(ctx context.Context, db *sqlx.DB, config env.EnvApp) Starter {
	gin.SetMode(config.GIN_MODE)
	server := &AppServer{
		gin.Default(),
		ctx,
		db,
		config,
	}

	server.configure()
	server.routes()
	return server
}

func (s *AppServer) Start() {
	fmt.Printf("Server running on port %s\n", s.config.PORT)
	s.Engine.Run(":" + s.config.PORT)
}
