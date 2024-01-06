package loadBalancer

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	HttpServerPort = 8080
	MaxRetries     = 3
)

func main() {
	var serverList string

	flag.StringVar(&serverList, "servers", "", "Load balanced servers. Use commas to separate")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Load balancer should have at least 1 server.")
	}

	// parse server list
	servers := strings.Split(serverList, ",")
	for _, s := range servers {
		serverUrl, err := url.Parse(s)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())

			retries := GetRetryFromContext(request)
			if retries < MaxRetries {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this server as down
			serverPool.ChangeServerStatus(serverUrl, false)

			// if the same request routing for few attempts with different servers, increase the count
			attempts := GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)

			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			loadBalancer(writer, request.WithContext(ctx))
		}

		serverPool.AddBackend(0,
			&Backend{
				URL:          serverUrl.String(),
				Alive:        true,
				ReverseProxy: proxy,
			})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", HttpServerPort),
		Handler: http.HandlerFunc(loadBalancer),
	}

	go healthCheck()

	log.Printf("Load Balancer is working at HttpServerPort: %d\n", HttpServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
