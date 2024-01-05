package proxyServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Address      string // Address contains URI of this particular server
	ServerConfig ServerConfig
	Router       *gin.Engine
}

type ServerConfig struct {
	Port string
}

// First of all, I need to set up server where I want to forward requests to.
// I will use http.Server struct for that purpose.
func initiateServer() *Server {
	s := &Server{
		Address: fmt.Sprintf("http://localhost:%s", HttpServerAddress),
		ServerConfig: ServerConfig{
			Port: HttpServerAddress,
		},
		Router: gin.Default(),
	}

	s.configureRoutes()

	return s
}

func (s *Server) configureRoutes() {
	s.Router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
}
