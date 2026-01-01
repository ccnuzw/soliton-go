package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// Import your modules here:
	// userapp "github.com/soliton-go/test-project/internal/application/user"
	// "github.com/soliton-go/test-project/internal/interfaces/http"
)

func main() {
	fx.New(
		// Database
		fx.Provide(NewDB),

		// Modules - uncomment after generating domains:
		// userapp.Module,

		// HTTP Handlers - uncomment after generating domains:
		// fx.Provide(http.NewUserHandler),

		// Start server
		fx.Invoke(StartServer),
	).Run()
}

// NewDB creates a new GORM database connection.
func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("✅ Database connected")
	return db
}

// StartServer starts the HTTP server.
func StartServer(db *gorm.DB) {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register routes - uncomment after generating domains:
	// userHandler.RegisterRoutes(r)

	log.Println("🚀 Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
