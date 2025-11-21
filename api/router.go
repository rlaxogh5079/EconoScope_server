package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rlaxogh5079/EconoScope/api/middleware"
	"github.com/rlaxogh5079/EconoScope/pkg/response"
)

func SetUpRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.GinLogger())
	
	r = middleware.SetCorsMiddleware(r)

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/ping", func(c* gin.Context) {
			response.Success(c, http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	return r
}