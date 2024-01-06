package proxyServer

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ProxyServer struct represents a proxy server
// that listens on a specific address and forwards
// requests to a forwarding handler
type ProxyServer struct {
	ForwardToServerAddress string // ForwardToServerAddress is forwarding address. IP/URI of server where this, proxy server, needs transport requests to.
	Config                 ProxyServerConfig
	Router                 *gin.Engine // Router is an alias for the proxy server
	CertFile               string      // Path to the .p12 certificate file
	CertPassword           string      // Password for the .p12 certificate
	DisableCertVerify      bool        // Whether to disable certificate verification for development
}

type ProxyServerConfig struct {
	Port string
}

// NewProxyServer creates a new proxyServer struct
// that will listen on a specific address and forwards
// requests to a forwarding handler
//
// Proxy server runs on port 8080, and will forward requests to Server's port (8081)
func NewProxyServer(forwardToServerPort string, certFile string, certPassword string, disableCertVerify bool) (*ProxyServer, error) {
	proxyServerConfig := ProxyServerConfig{
		Port: ProxyServerHttpAddress,
	}

	router := gin.Default()
	router.GET("/", forwardingHandler)

	return &ProxyServer{
		ForwardToServerAddress: forwardToServerPort,
		Config:                 proxyServerConfig,
		Router:                 router,
		CertFile:               certFile,
		CertPassword:           certPassword,
		DisableCertVerify:      disableCertVerify,
	}, nil
}

// Listen serves on a specific address and forwards
// requests to a forwarding handler
func (p *ProxyServer) Listen() error {
	// Load the .p12 certificate and configure the HTTP client
	clientCert, err := tls.LoadX509KeyPair(p.CertFile, p.CertPassword)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
	}

	if p.DisableCertVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: transport}

	http.DefaultClient = client

	// Start the proxy server
	if err := p.Router.Run(fmt.Sprintf(":%s", p.Config.Port)); err != nil {
		ErrorLogger.Printf("Error running server on port 8080: %v\n", err)
		return err
	}

	return nil
}

// forwardingHandler forwards requests to the given address
func forwardingHandler(c *gin.Context) {
	// firstly, create a new request and copy all the data from the original request

	// create a new request
	proxyServerRequest := &http.Request{}
	*proxyServerRequest = *c.Request
	proxyServerRequest.RequestURI = "http://localhost:" + fmt.Sprintf("%s", HttpServerAddress)

	response, err := http.DefaultTransport.RoundTrip(proxyServerRequest)
	if err != nil {
		ErrorLogger.Printf("Error forwarding request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// copy original request headers
	for header, values := range c.Request.Header {
		for _, value := range values {
			c.Writer.Header().Add(header, value)
		}
	}

	c.Writer.WriteHeader(response.StatusCode)

	_, copyErr := io.Copy(c.Writer, response.Body)
	if copyErr != nil {
		ErrorLogger.Printf("Error copying response: %v\n", copyErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": copyErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "response",
	})
}
