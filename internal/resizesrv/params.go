package resizesrv

import (
	"fmt"
	"strconv"
)

func validateParams(width, height, url string) (*ImageParams, error) {
	w, err := validateSize(width, "width")
	if err != nil {
		return nil, err
	}
	h, err := validateSize(height, "height")
	if err != nil {
		return nil, err
	}
	return &ImageParams{Width: w, Height: h, URL: url}, nil
}

func validateSize(size string, t string) (int, error) {
	w, err := strconv.Atoi(size)
	if err != nil {
		return 0, fmt.Errorf("%s must be int", t)
	}
	if w <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", t)
	}
	return w, nil
}
