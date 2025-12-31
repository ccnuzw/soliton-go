package web

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

// Server is a wrapper around Gin engine.
type Server struct {
	engine *gin.Engine
}

// NewServer creates a new Server instance.
func NewServer() *Server {
	r := gin.Default()
	return &Server{engine: r}
}

// RegisterGraphQL registers a GraphQL handler at the given path.
func (s *Server) RegisterGraphQL(path string, server *handler.Server) {
	s.engine.POST(path, func(c *gin.Context) {
		server.ServeHTTP(c.Writer, c.Request)
	})
	s.engine.GET(path+"/playground", func(c *gin.Context) {
		playground.Handler("GraphQL", path).ServeHTTP(c.Writer, c.Request)
	})
}

// Run starts the HTTP server.
func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
