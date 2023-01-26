package imagesrv

import (
	"bytes"
	"context"
	"image"
	"net/http"
)

func New(cancelContext context.Context) *ImageSrv {
	return &ImageSrv{
		cancelContext: cancelContext,
		httpClient:    http.DefaultClient,
	}
}

func (is *ImageSrv) ExtractParams(path string) (*ImageParams, error) {
	return is.extractParams(path)
}

func (is *ImageSrv) GetImg(url string) (*image.Image, error) {
	file, err := is.downloadFile(url)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	return &img, err
}
