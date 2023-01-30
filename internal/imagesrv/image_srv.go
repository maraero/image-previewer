package imagesrv

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/maraero/image-previewer/internal/cache"
	"github.com/maraero/image-previewer/internal/logger"
)

func New(cancelContext context.Context, cache cache.Cache, logger logger.Logger) *ImageSrv {
	return &ImageSrv{
		cache:         cache,
		cancelContext: cancelContext,
		httpClient:    http.DefaultClient,
		logger:        logger,
	}
}

func (is *ImageSrv) GetResizedImg(params string) ([]byte, error) {
	cacheKey := getCacheKey(params)
	cachedImg, exists := is.cache.Get(cacheKey)
	if exists {
		is.logger.Info("get image from cache by key ", cacheKey)
		return cachedImg, nil
	}

	imgParams, err := extractParams(params)
	if err != nil {
		is.logger.Error(err)
		return nil, err
	}

	img, err := is.getImg(imgParams.URL)
	if err != nil {
		is.logger.Error(err)
		return nil, err
	}

	rszdImg := is.resizeImage(img, imgParams.Width, imgParams.Height)

	imgBytes, err := is.encodeImageToBytes(&rszdImg)
	if err != nil {
		is.logger.Error(err)
		return nil, ErrEncodingToBytes
	}

	err = is.cache.Set(cacheKey, imgBytes)
	if err != nil {
		is.logger.Error("%s: %w", ErrCacheSet, err)
		return nil, ErrCacheSet
	}

	return imgBytes, nil
}

func (is *ImageSrv) getImg(url string) (*image.Image, error) {
	file, err := is.downloadFile(url)
	if err != nil {
		return nil, err
	}

	jpeg := is.isFileJPEG(file[0:3])
	if !jpeg {
		return nil, fmt.Errorf("%s %w", url, ErrFileIsNotJPEG)
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, ErrCanNotDecodeJPEG
	}

	return &img, nil
}

func (is *ImageSrv) encodeImageToBytes(img *image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *img, nil)
	if err != nil {
		return []byte{}, ErrEncodingToBytes
	}
	return buf.Bytes(), nil
}

func (is *ImageSrv) resizeImage(img *image.Image, width, height int) image.Image {
	return imaging.Thumbnail(*img, width, height, imaging.Lanczos)
}

func (is *ImageSrv) downloadFile(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(is.cancelContext, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCanNotBuildRequest, err)
	}

	resp, err := is.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCanNotMakeRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w from %s", ErrCanNotDownloadFile, url)
	}

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCanNotReadResponseBody, err)
	}

	return image, nil
}

func (is *ImageSrv) isFileJPEG(firstFileBytes []byte) bool {
	for i, bt := range firstFileBytes {
		if bt != JPEGMagicNumber[i] {
			return false
		}
	}
	return true
}

func getCacheKey(params string) string {
	return cacheKeyRegexp.ReplaceAllString(params, "_")
}
