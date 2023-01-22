package imagesrv

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func New() *ImageSrv {
	return &ImageSrv{}
}

func (rs *ImageSrv) ExtractParams(path string) (*ImageParams, error) {
	p := strings.Split(path, "/")
	if len(p) < 3 {
		return nil, ErrTooFewParams
	}
	width := p[0]
	height := p[1]
	url := strings.Join(p[2:], "/")

	return validateParams(width, height, url)
}

func (rs *ImageSrv) DownloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("can not download file from %s", url)
	}
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return image, nil
}
