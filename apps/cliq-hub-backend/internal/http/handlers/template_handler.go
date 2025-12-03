package handlers

import (
	"net/http"

	"cliq-hub-backend/internal/auth"
	"cliq-hub-backend/internal/db"
	"cliq-hub-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct{}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) List(c *gin.Context) {
	var templates []models.Template
	if result := db.DB.Preload("Author").Find(&templates); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var template models.Template
	if result := db.DB.Preload("Author").First(&template, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *TemplateHandler) Create(c *gin.Context) {
	// Simple auth middleware check
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Content     string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := models.Template{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		AuthorID:    userID.(uint),
	}

	if result := db.DB.Create(&template); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// AuthMiddleware extracts the user ID from the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// Handle "Bearer <token>"
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims, err := auth.ValidateToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
