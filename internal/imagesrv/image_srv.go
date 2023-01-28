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
	"github.com/maraero/image-previewer/internal/logger"
)

func New(cancelContext context.Context, logger logger.Logger) *ImageSrv {
	return &ImageSrv{
		cancelContext: cancelContext,
		httpClient:    http.DefaultClient,
		logger:        logger,
	}
}

func (is *ImageSrv) GetResizedImg(params string) ([]byte, error) {
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

	return imgBytes, nil
}

func (is *ImageSrv) getImg(url string) (*image.Image, error) {
	file, err := is.downloadFile(url)
	if err != nil {
		is.logger.Error(err)
		return nil, ErrFileDownload
	}

	jpeg := is.isFileJPEG(file[0:3])
	if !jpeg {
		is.logger.Error("%s is not jpeg", url)
		return nil, ErrFileIsNotJPEG
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		is.logger.Error(err)
		return nil, ErrCanNotDecodeJPEG
	}

	return &img, nil
}

func (is *ImageSrv) encodeImageToBytes(img *image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *img, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("can not encode image: %w", err)
	}
	return buf.Bytes(), nil
}

func (is *ImageSrv) resizeImage(img *image.Image, width, height int) image.Image {
	return imaging.Thumbnail(*img, width, height, imaging.Lanczos)
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

func (is *ImageSrv) isFileJPEG(firstFileBytes []byte) bool {
	for i, bt := range firstFileBytes {
		if bt != JPEGMagicNumber[i] {
			return false
		}
	}
	return true
}
