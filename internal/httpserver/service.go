package httpserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/maraero/image-previewer/internal/app"
	"github.com/maraero/image-previewer/internal/imagesrv"
	"github.com/maraero/image-previewer/internal/logger"
)

func handleFill(app *app.App, l logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sep := "/fill/"
		p := strings.Split(r.URL.Path, sep)
		params := strings.Join(p[1:], sep)
		reqHeaders := r.Header.Clone()
		image, err := app.ImageSrv.GetResizedImg(params, reqHeaders)
		var paramValidationError *imagesrv.ParamValidationError

		// Bad Request
		if err != nil && (errors.As(err, &paramValidationError) || errors.Is(err, imagesrv.ErrIsNotJPEG)) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		// Bad Gateway
		if err != nil && errors.Is(err, imagesrv.ErrCanNotDownloadFile) {
			w.WriteHeader(http.StatusBadGateway)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		// Internal Server Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(image); err != nil {
			l.Error("http write error: %w", err)
		}
	}
}
