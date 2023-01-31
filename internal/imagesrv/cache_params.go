package imagesrv

import (
	"net/http"
	"time"

	"github.com/maraero/image-previewer/internal/cache"
	"github.com/marcw/cachecontrol"
)

func collectCacheParams(respHeaders http.Header) cache.Params {
	var expiresAt time.Time
	var lastModifiedAt time.Time

	expires := respHeaders.Get("expires")
	expiresAt, _ = time.Parse(http.TimeFormat, expires)

	// cache-control has a higher priority than expires
	cacheControl := respHeaders.Get("cache-control")
	if cacheControl != "" {
		cacheControlObj := cachecontrol.Parse(cacheControl)
		var cacheDuration time.Duration

		if public, _ := cacheControlObj.Private(); public {
			cacheDuration = cacheControlObj.MaxAge()
		}

		if cacheDuration != 0 {
			expiresAt = time.Now().Add(cacheDuration)
		}
	}

	etag := respHeaders.Get("etag")

	lastModified := respHeaders.Get("last-modified")
	lastModifiedAt, _ = time.Parse(http.TimeFormat, lastModified)

	return cache.Params{
		ExpiresAt:      expiresAt,
		Etag:           etag,
		LastModifiedAt: lastModifiedAt,
	}
}
