package router

import (
	"cliq-hub-backend/internal/db"
	"cliq-hub-backend/internal/http/handlers"
	"cliq-hub-backend/internal/llm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New(client llm.Client, debugMode bool) *gin.Engine {
	// Initialize DB
	db.Init("") // Use default path

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	genHandler := handlers.NewGenerateHandler(client, debugMode)
	authHandler := handlers.NewAuthHandler()
	templateHandler := handlers.NewTemplateHandler()

	v1 := r.Group("/v1")

	// Auth
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Templates
	tm := v1.Group("/templates")
	tm.POST("/generate", genHandler.Handle) // Existing
	tm.GET("", templateHandler.List)
	tm.GET("/:id", templateHandler.Get)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/templates", templateHandler.Create)
	}

	return r
}
