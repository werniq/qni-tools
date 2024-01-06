package loadBalancer

import (
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
)

// ServerPool holds information about reachable backends
type ServerPool struct {
	backends *BSTBalancer
	current  uint64
}

// AddBackend to the server pool tree
func (s *ServerPool) AddBackend(key int, backend *Backend) {
	s.backends.Insert(key, backend)
}

// NextIndex atomically increase the counter and return an index
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(s.backends.LastKey))
}

// ChangeServerStatus changes a status of a backend
func (s *ServerPool) ChangeServerStatus(backendUrl *url.URL, alive bool) {
	url := backendUrl.String()

	for i := 0; i < s.backends.Key; i++ {
		if s.backends.Search(i).Val.URL == url {
			s.backends.Val.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	// start from next and move a full cycle
	l := s.backends.Key + next

	for i := next; i < l; i++ {
		// take an index by mod operation
		idx := i % s.backends.Max()
		if s.backends.Search(idx).Val.IsAlive() {
			// if we have an alive backend, use it and store if it's not the original one
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}

			return s.backends.Search(idx).Val
		}
	}

	return nil
}

// InOrderHealthCheck pings the backends and update the status
func (s *ServerPool) InOrderHealthCheck(b *BSTBalancer) {
	for b != nil {
		s.InOrderHealthCheck(b.Left)

		status := "up"

		alive := isBackendAlive(b.Val.URL)
		b.Val.SetAlive(alive)

		if !alive {
			status = "down"
		}

		log.Printf("%d | %s | is [%s]\n", b.Key, b.Val.URL, status)

		s.InOrderHealthCheck(b.Right)
	}
}

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// GetRetryFromContext returns the retries for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}
