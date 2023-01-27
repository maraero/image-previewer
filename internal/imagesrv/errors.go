package imagesrv

import "errors"

type ParamValidationError struct {
	desc string
}

func (e *ParamValidationError) Error() string {
	return e.desc
}

var (
	ErrFileDownload    = errors.New("can not download file")
	ErrFileIsNotJpeg   = errors.New("is not jpeg")
	ErrEncodingToBytes = errors.New("can not encode image to bytes")
)

const ErrorTooFewParams = "too few params. Pass width, height and url in format /fill/:width/:height/:url"
