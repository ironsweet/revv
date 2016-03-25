package main

import (
	"net/http"
	"time"
)

func main() {
	now := time.Now()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(now.String()))
	})
	http.ListenAndServe(":8080", nil)
}
