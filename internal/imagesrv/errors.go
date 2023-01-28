package imagesrv

import "errors"

type ParamValidationError struct {
	desc string
}

func (e *ParamValidationError) Error() string {
	return e.desc
}

var (
	ErrCanNotBuildRequest     = errors.New("can not build request")
	ErrCanNotDecodeJPEG       = errors.New("can not decode jpeg")
	ErrCanNotDownloadFile     = errors.New("can not download file")
	ErrCanNotMakeRequest      = errors.New("can not make request")
	ErrCanNotReadResponseBody = errors.New("can not read response body")
	ErrEncodingToBytes        = errors.New("can not encode image to bytes")
	ErrFileIsNotJPEG          = errors.New("is not jpeg")
)

const ErrorTooFewParams = "too few params. Pass width, height and url in format /fill/:width/:height/:url"
