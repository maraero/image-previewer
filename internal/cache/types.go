package cache

import "time"

type Params struct {
	ExpiresAt      time.Time
	Etag           string
	LastModifiedAt time.Time
}
