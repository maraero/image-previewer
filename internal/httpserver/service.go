package httpserver

import (
	"net/http"

	"github.com/maraero/image-previewer/internal/logger"
)

func handleHello(l logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello"))
		if err != nil {
			l.Error("http write error: %w", err)
		}
	}
}
