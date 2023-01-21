package httpserver

import (
	"fmt"
	"net/http"
)

func handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello"))
		if err != nil {
			fmt.Println("http write error: %w", err)
		}
	}
}
