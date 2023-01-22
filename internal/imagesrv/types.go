package imagesrv

import (
	"context"
	"net/http"
)

type ImageSrv struct {
	cancelContext context.Context
	httpClient    *http.Client
}

type ImageParams struct {
	Width  int
	Height int
	URL    string
}
