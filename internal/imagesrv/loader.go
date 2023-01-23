package imagesrv

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (is *ImageSrv) extractParams(path string) (*ImageParams, error) {
	p := strings.Split(path, "/")
	if len(p) < 3 {
		return nil, ErrTooFewParams
	}
	width := p[0]
	height := p[1]
	url := strings.Join(p[2:], "/")

	return validateParams(width, height, url)
}

func (is *ImageSrv) downloadFile(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(is.cancelContext, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("can not make request")
	}
	resp, err := is.httpClient.Do(req)
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
