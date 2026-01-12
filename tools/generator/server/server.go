package server

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/soliton-go/tools/server/handlers"
)

//go:embed all:static
var staticFS embed.FS

// Start starts the web server.
func Start(host string, port int) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
	{
		// Project endpoints
		api.POST("/projects/init", handlers.InitProject)
		api.POST("/projects/init/preview", handlers.PreviewInitProject)
		api.POST("/projects/switch", handlers.SwitchProject)
		api.GET("/projects/list", handlers.ListProjects)
		api.POST("/projects/tidy", handlers.RunGoModTidy)
		api.POST("/projects/migrate", handlers.RunMigration)

		// Domain endpoints
		api.POST("/domains", handlers.GenerateDomain)
		api.POST("/domains/preview", handlers.PreviewDomain)
		api.GET("/domains/list", handlers.ListDomains)
		api.GET("/domains/:name", handlers.GetDomainDetail)
		api.DELETE("/domains/:name", handlers.DeleteDomain)
		api.GET("/field-types", handlers.GetFieldTypes)

		// Service endpoints
		api.POST("/services", handlers.GenerateService)
		api.POST("/services/preview", handlers.PreviewService)
		api.GET("/services/list", handlers.ListServices)
		api.GET("/services/detect/:name", handlers.DetectServiceType)
		api.GET("/services/:name", handlers.GetServiceDetail)
		api.DELETE("/services/:name", handlers.DeleteService)

		// DDD endpoints
		api.POST("/ddd/valueobjects", handlers.GenerateValueObject)
		api.POST("/ddd/valueobjects/preview", handlers.PreviewValueObject)
		api.POST("/ddd/specs", handlers.GenerateSpecification)
		api.POST("/ddd/specs/preview", handlers.PreviewSpecification)
		api.POST("/ddd/policies", handlers.GeneratePolicy)
		api.POST("/ddd/policies/preview", handlers.PreviewPolicy)
		api.POST("/ddd/events", handlers.GenerateEvent)
		api.POST("/ddd/events/preview", handlers.PreviewEvent)
		api.POST("/ddd/event-handlers", handlers.GenerateEventHandler)
		api.POST("/ddd/event-handlers/preview", handlers.PreviewEventHandler)
		api.GET("/ddd/list", handlers.ListDDD)
		api.GET("/ddd/detail", handlers.GetDDDDetail)
		api.GET("/ddd/source", handlers.GetDDDSource)
		api.POST("/ddd/delete", handlers.DeleteDDD)
		api.POST("/ddd/rename", handlers.RenameDDD)

		// Utility endpoints
		api.GET("/layout", handlers.GetProjectLayout)
	}

	// Static files - try embedded first, then fall back to file system
	staticHandler := createStaticHandler()
	r.NoRoute(staticHandler)

	// Start server
	addr := fmt.Sprintf("%s:%d", host, port)
	if err := r.Run(addr); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".html":
		return "text/html; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

func createStaticHandler() gin.HandlerFunc {
	// Try to use embedded static files
	subFS, err := fs.Sub(staticFS, "static")
	if err == nil {
		// Check if index.html exists in embedded FS
		if _, err := subFS.Open("index.html"); err == nil {
			// Use custom handler for SPA routing
			return func(c *gin.Context) {
				path := c.Request.URL.Path

				// Remove leading slash for fs.FS
				cleanPath := path
				if len(cleanPath) > 0 && cleanPath[0] == '/' {
					cleanPath = cleanPath[1:]
				}

				// Try to open the requested file
				f, err := subFS.Open(cleanPath)
				if err == nil {
					defer f.Close()
					stat, err := f.Stat()
					if err == nil && !stat.IsDir() {
						// File exists, serve it directly with correct content type
						contentType := getContentType(cleanPath)
						c.Header("Content-Type", contentType)
						http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), f.(io.ReadSeeker))
						return
					}
				}

				// File not found, serve index.html for SPA routing
				indexFile, err := subFS.Open("index.html")
				if err == nil {
					defer indexFile.Close()
					stat, _ := indexFile.Stat()
					c.Header("Content-Type", "text/html; charset=utf-8")
					http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
				}
			}
		}
	}

	// Fall back to serving from file system (for development)
	return func(c *gin.Context) {
		// Try to find web/dist directory
		possiblePaths := []string{
			"web/dist",
			"tools/generator/web/dist",
			filepath.Join(os.Getenv("HOME"), "Progame/soliton-go/tools/generator/web/dist"),
		}

		for _, basePath := range possiblePaths {
			indexPath := filepath.Join(basePath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				// Serve static files from this directory
				requestedPath := c.Request.URL.Path
				filePath := filepath.Join(basePath, requestedPath)

				// If file exists, serve it
				if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
					c.File(filePath)
					return
				}

				// Otherwise serve index.html for SPA routing
				c.File(indexPath)
				return
			}
		}

		// No static files found, return a simple HTML page
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `<!DOCTYPE html>
<html>
<head>
    <title>Soliton-Gen Web GUI</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
               max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .api-list { background: #f5f5f5; padding: 20px; border-radius: 8px; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 4px; }
    </style>
</head>
<body>
    <h1>🚀 Soliton-Gen Web GUI</h1>
    <p>The API server is running. The frontend is not yet built.</p>
    
    <div class="api-list">
        <h3>Available API Endpoints:</h3>
        <ul>
            <li><code>POST /api/projects/init</code> - Initialize a new project</li>
            <li><code>POST /api/projects/init/preview</code> - Preview project initialization</li>
            <li><code>POST /api/projects/migrate</code> - Run migrations with logs</li>
            <li><code>POST /api/domains</code> - Generate a domain module</li>
            <li><code>POST /api/domains/preview</code> - Preview domain generation</li>
            <li><code>GET /api/field-types</code> - Get available field types</li>
            <li><code>POST /api/services</code> - Generate a service</li>
            <li><code>POST /api/services/preview</code> - Preview service generation</li>
            <li><code>POST /api/ddd/valueobjects</code> - Generate a value object</li>
            <li><code>POST /api/ddd/specs</code> - Generate a specification</li>
            <li><code>POST /api/ddd/policies</code> - Generate a policy</li>
            <li><code>POST /api/ddd/events</code> - Generate a domain event</li>
            <li><code>POST /api/ddd/event-handlers</code> - Generate an event handler</li>
            <li><code>GET /api/layout</code> - Get current project layout</li>
        </ul>
    </div>
    
    <p style="margin-top: 20px; color: #666;">
        To build the frontend, run: <code>cd web && npm install && npm run build</code>
    </p>
</body>
</html>`)
	}
}
