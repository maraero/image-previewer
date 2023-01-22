package httpserver

import (
	"net/http"
	"strings"

	"github.com/maraero/image-previewer/internal/app"
	"github.com/maraero/image-previewer/internal/logger"
)

func handleFill(app *app.App, l logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split(r.URL.Path, "/")
		params := strings.Join(p[2:], "/")
		ip, err := app.ResizeSrv.ExtractParams(params)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		l.Info(ip.Width, ip.Height, ip.URL)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Success"))
		if err != nil {
			l.Error("http write error: %w", err)
		}
	}
}
