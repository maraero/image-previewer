package imagesrv

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"
)

func (is *ImageSrv) downloadJPEG(url string, headers http.Header) (*image.Image, error) {
	req, err := http.NewRequestWithContext(is.cancelContext, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCanNotBuildRequest, err)
	}

	for h := range headers {
		req.Header.Add(h, headers.Get(h))
	}

	resp, err := is.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCanNotMakeRequest, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w from %s", ErrCanNotDownloadFile, url)
	}

	respHeaders := resp.Header.Clone()
	if !is.canBeJPEG(respHeaders) {
		return nil, ErrIsNotJPEG
	}

	cacheParams := collectCacheParams(respHeaders)
	fmt.Println(cacheParams) // TODO: Delete

	defer resp.Body.Close()
	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCanNotReadResponseBody, err)
	}

	jpeg := is.isFileJPEG(file[0:3])
	if !jpeg {
		return nil, ErrIsNotJPEG
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, ErrCanNotDecodeJPEG
	}

	return &img, nil
}

func (is *ImageSrv) canBeJPEG(respHeaders http.Header) bool {
	ct := respHeaders.Get("content-type")
	if ct == "" || strings.Contains(ct, "image/jpeg") {
		return true
	}
	return false
}

func (is *ImageSrv) isFileJPEG(firstFileBytes []byte) bool {
	for i, bt := range firstFileBytes {
		if bt != JPEGMagicNumber[i] {
			return false
		}
	}
	return true
}
