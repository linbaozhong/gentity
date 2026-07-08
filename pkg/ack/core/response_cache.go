// Copyright © 2023 Linbaozhong. All rights reserved.
// ...license same as before...

package core

import (
	"context"
	"net/http"
	"time"

	"github.com/linbaozhong/gentity/pkg/cachego/mmap"
)

type cacheKey struct {
	Key      string
	Duration time.Duration
}

var respCache = mmap.New(mmap.WithExpired(time.Second * 30))

// ReadCache checks if a cached response exists for the current GET request.
// If found, writes it directly and returns true (handler should return immediately).
func ReadCache(ctx Context, lefetime ...time.Duration) bool {
	u := ctx.Request().URL
	vals := u.Values
	if vals == nil {
		vals = make(map[string][]string)
	}
	vals.Set("_t", ctx.GetHeader(authorizationHdr))
	key := u.Path + "?" + vals.Encode()

	duration := time.Second * 30
	if len(lefetime) > 0 {
		duration = lefetime[0]
	}
	ctx.Values().Set(hasCacheKey, cacheKey{Key: key, Duration: duration})

	buf, e := respCache.Fetch(ctx, key)
	if e != nil || len(buf) == 0 {
		return false
	}

	ctx.StopWithStatus(http.StatusOK)
	ctx.SetContentType("application/json")
	_, e = ctx.ResponseWriter().Write(buf)
	return e == nil
}

func setCache(ctx context.Context, key cacheKey, val []byte) {
	respCache.Save(ctx, key.Key, val, key.Duration)
}
