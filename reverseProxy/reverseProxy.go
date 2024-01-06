package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	ErrorLogger = log.New(os.Stdout, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger  = log.New(os.Stdout, "INFO: \t", log.Ldate|log.Ltime|log.Lshortfile)
)

// launching reverse proxy and main server
func main() {
	originServerUri, err := url.Parse("https://localhost:8081")
	if err != nil {
		ErrorLogger.Printf("Error parsing uri: %v\n", err)
		return
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		InfoLogger.Printf("[reverse proxy server] received request at: %s\n", time.Now())

		// change the request to point to the origin server
		r.Host = originServerUri.Host
		r.URL.Host = originServerUri.Host
		r.Host = originServerUri.Host
		r.URL.Scheme = originServerUri.Scheme
		r.RequestURI = ""

		originServerResponse, err := http.DefaultClient.Do(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}
		defer originServerResponse.Body.Close()

		InfoLogger.Printf("[origin server] response status: %s\n", originServerResponse.Status)

		rw.WriteHeader(http.StatusOK)
		_, err = io.Copy(rw, originServerResponse.Body)
		if err != nil {
			ErrorLogger.Printf("Error copying response: %v\n", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}
