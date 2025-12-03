package db

import (
	"log"
	"path/filepath"
	"runtime"

	"cliq-hub-backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dbPath string) {
	if dbPath == "" {
		// Default to a file in the project root if not specified
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		// Go up to apps/cliq-hub-backend
		root := filepath.Join(basepath, "../../..")
		dbPath = filepath.Join(root, "cliqhub.db")
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Auto Migrate
	err = DB.AutoMigrate(&models.User{}, &models.Template{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("Database initialized at", dbPath)
}
