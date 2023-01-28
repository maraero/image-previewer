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

		image, err := app.ImageSrv.GetResizedImg(params)
		var paramValidationError *imagesrv.ParamValidationError

		if err != nil && (errors.As(err, &paramValidationError) || errors.Is(err, imagesrv.ErrFileIsNotJPEG)) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		if err != nil && errors.Is(err, imagesrv.ErrFileDownload) {
			w.WriteHeader(http.StatusBadGateway)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				l.Error("http write error: %w", err)
			}
			return
		}

		if err != nil && errors.Is(err, imagesrv.ErrEncodingToBytes) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(image); err != nil {
			l.Error("http write error: %w", err)
		}
	}
}
