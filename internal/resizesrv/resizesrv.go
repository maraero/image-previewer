package resizesrv

import (
	"strings"
)

func New() *ResizeSrv {
	return &ResizeSrv{}
}

func (rs *ResizeSrv) ExtractParams(path string) (*ImageParams, error) {
	p := strings.Split(path, "/")
	if len(p) < 3 {
		return nil, ErrTooFewParams
	}
	width := p[0]
	height := p[1]
	url := strings.Join(p[2:], "/")

	return validateParams(width, height, url)
}
