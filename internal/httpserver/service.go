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
		ip, err := app.ImageSrv.ExtractParams(params)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		image, err := app.ImageSrv.DownloadImage(ip.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(image); err != nil {
			l.Error("http write error: %w", err)
		}
	}
}
