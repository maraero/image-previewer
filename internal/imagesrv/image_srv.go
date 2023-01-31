package imagesrv

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
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

func (is *ImageSrv) GetResizedImg(params string, reqHeaders http.Header) ([]byte, error) {
	cacheKey := getCacheKey(params)
	cachedImg, exists, _ := is.cache.Get(cacheKey)
	if exists {
		is.logger.Info("get image from cache by key ", cacheKey)
		return cachedImg, nil
	}

	imgParams, err := extractParams(params)
	if err != nil {
		is.logger.Error(err)
		return nil, err
	}

	img, err := is.downloadJPEG(imgParams.URL, reqHeaders)
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

func (is *ImageSrv) encodeImageToBytes(img *image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *img, nil)
	if err != nil {
		return []byte{}, ErrEncodingToBytes
	}
	return buf.Bytes(), nil
}

func (is *ImageSrv) isJPEG(respHeaders http.Header) bool {
	ct, ok := respHeaders["content-type"]
	if !ok {
		return true
	}
	for _, v := range ct {
		if v == "image/jpeg" {
			return true
		}
	}
	return false
}

func (is *ImageSrv) resizeImage(img *image.Image, width, height int) image.Image {
	return imaging.Fill(*img, width, height, imaging.Center, imaging.Lanczos)
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
