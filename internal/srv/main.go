package main

import (
	"fmt"
	"net/http"
)

func main() {
	println("serving on port :9090")
	err := http.ListenAndServe(":9090", http.FileServer(http.Dir("testdata")))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
