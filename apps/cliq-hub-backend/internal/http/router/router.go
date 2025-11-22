package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"cliq-hub-backend/internal/http/handlers"
	"cliq-hub-backend/internal/llm"
)

func New(client llm.Client, debugMode bool) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
	}))
	h := handlers.NewGenerateHandler(client, debugMode)

	v1 := r.Group("/v1")
	tm := v1.Group("/templates")
	tm.POST("/generate", h.Handle)
	return r
}
