package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cliq-hub-backend/internal/db"
	"cliq-hub-backend/internal/http/handlers"
	"cliq-hub-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Init DB
	db.Init(os.TempDir() + "/test.db")
	db.DB.AutoMigrate(&models.User{}, &models.Template{})

	authHandler := handlers.NewAuthHandler()
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	templateHandler := handlers.NewTemplateHandler()
	r.GET("/templates", templateHandler.List)
	r.GET("/templates/:id", templateHandler.Get)

	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	protected.POST("/templates", templateHandler.Create)

	return r
}

func TestAuthAndTemplateFlow(t *testing.T) {
	r := setupRouter()

	// 1. Register
	registerPayload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Login
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)

	// 3. Create Template
	templatePayload := map[string]string{
		"title":   "Test Template",
		"content": "test: content",
	}
	body, _ = json.Marshal(templatePayload)
	req, _ = http.NewRequest("POST", "/templates", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 4. List Templates
	req, _ = http.NewRequest("GET", "/templates", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var templates []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &templates)
	assert.Len(t, templates, 1)
	assert.Equal(t, "Test Template", templates[0]["title"])
}
