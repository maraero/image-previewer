package imagesrv

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/disintegration/imaging"
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

func (is *ImageSrv) DownloadImage(url string) (*image.Image, error) {
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

func (is *ImageSrv) EncodeImageToBytes(img *image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *img, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("can not encode image: %w", err)
	}
	return buf.Bytes(), nil
}

func (is *ImageSrv) ResizeImage(img *image.Image, width, height int) image.Image {
	return imaging.Thumbnail(*img, width, height, imaging.Lanczos)
}
