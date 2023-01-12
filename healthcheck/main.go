package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/health", os.Getenv("URL_COLLECTOR_PORT")))
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}

	resp.Body.Close()
}
