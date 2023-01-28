package imagesrv

import (
	"context"
	"net/http"

	"github.com/maraero/image-previewer/internal/logger"
)

type ImageService interface {
	GetResizedImg(params string) ([]byte, error)
}

type ImageSrv struct {
	cancelContext context.Context
	httpClient    *http.Client
	logger        logger.Logger
}

type ImageParams struct {
	Width  int
	Height int
	URL    string
}
