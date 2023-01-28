package imagesrv

import "regexp"

// JPEGMagicNumber is defined here: https://en.wikipedia.org/wiki/JPEG.
var JPEGMagicNumber = [3]byte{0xff, 0xd8, 0xff}

var cacheKeyRegexp = regexp.MustCompile("[^a-zA-Z0-9_]")
