package imagesrv

import (
	"strconv"
	"strings"
)

func extractParams(path string) (*ImageParams, error) {
	p := strings.Split(path, "/")
	if len(p) < 3 {
		return nil, &ParamValidationError{desc: ErrorTooFewParams}
	}
	width := p[0]
	height := p[1]
	url := getURLFromPath(path, width, height)
	return validateParams(width, height, url)
}

func getURLFromPath(path string, width string, height string) string {
	lpad := len(width) + len(height) + 2
	return path[lpad:]
}

func validateParams(width, height, url string) (*ImageParams, error) {
	w, err := validateSize(width, "width")
	if err != nil {
		return nil, err
	}
	h, err := validateSize(height, "height")
	if err != nil {
		return nil, err
	}
	return &ImageParams{Width: w, Height: h, URL: addSchemaToURLIfRequired(url)}, nil
}

func validateSize(size string, t string) (int, error) {
	w, err := strconv.Atoi(size)
	if err != nil {
		return 0, &ParamValidationError{desc: t + " must be int"}
	}
	if w <= 0 {
		return 0, &ParamValidationError{desc: t + " must be greater than zero"}
	}
	return w, nil
}

func addSchemaToURLIfRequired(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "http://" + url
}
