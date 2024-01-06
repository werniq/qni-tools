package loadBalancer

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

const (
	Attempts int = iota
	Retry
)

// Backend holds the data about a server
type Backend struct {
	Root         *BSTBalancer
	URL          string
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

// loadBalancer load balances the incoming request
func loadBalancer(w http.ResponseWriter, r *http.Request) {
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// isAlive checks whether a backend is Alive by establishing a TCP connection
func isBackendAlive(u string) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}

	defer conn.Close()
	return true
}

// healthCheck runs a routine for check status of the backends every 2 mins
func healthCheck() {
	t := time.NewTicker(time.Minute * 2)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.InOrderHealthCheck(serverPool.backends)
			log.Println("Health check completed")
		}
	}
}

var serverPool ServerPool
