package imagesrv

import "errors"

var ErrTooFewParams = errors.New("too few params. Pass width, height and url in format /fill/:width/:height/:url")
