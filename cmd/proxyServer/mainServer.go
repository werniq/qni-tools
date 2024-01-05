package proxyServer

import (
	"fmt"
	"log"
	"os"
)

var (
	ErrorLogger = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger  = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)
)

const (
	HttpServerAddress      = "8080"
	ProxyServerHttpAddress = "8081"
)

func main() {
	server := initiateServer()

	pathToCertFile := "path"
	passwordToCertFile := "super-secret-password"
	disableCertVerification := true // disable certificate verification for development

	proxyServer, err := NewProxyServer(server.ServerConfig.Port, pathToCertFile, passwordToCertFile, disableCertVerification)

	if err != nil {
		ErrorLogger.Printf(err.Error())
		return
	}

	proxyServer.Router.Run(fmt.Sprintf(":%s", proxyServer.Config.Port))

	server.Router.Run(fmt.Sprintf(":%s", server.ServerConfig.Port))

}
