package server

import (
	"support/internal/core/services"
	"support/internal/handlers"
	"support/internal/providers"
	"support/internal/repositories"

	"github.com/gin-gonic/gin"
)

func healthRoutes(ginServer *AppServer) {
	ginServer.Engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func statsRoutes(ginServer *AppServer) {
	ginServer.Engine.GET("/stats", func(c *gin.Context) {
		var routeStatsList []*handlers.RouteStats
		for _, routeStats := range handlers.GetStatsMap() {
			if routeStats.Path != "/stats" {
				routeStatsList = append(routeStatsList, routeStats)
			}
		}

		c.JSON(200, routeStatsList)
	})

}

func securityRoutes(appServer *AppServer) {
	group := appServer.Engine.Group("/security")
	repository := repositories.New()
	jwtService := providers.NewJwt(appServer.config)
	redisRepository := providers.NewRedis(appServer.ctx, appServer.config)
	service := services.New(repository, jwtService, redisRepository)
	cidiService := services.NewCidi(appServer.ctx, appServer.config)
	controller := handlers.New(service, cidiService)

	group.GET("temporal-user", controller.UsuarioTemporal)
	group.GET("login", controller.Login)
	group.POST("permissions", controller.Permisos)
	group.POST("menu", controller.Menu)
	group.GET("has-represented", controller.ObtenerRepresentado)
}

func (s *AppServer) routes() {
	healthRoutes(s)
	securityRoutes(s)
	statsRoutes(s)
}
