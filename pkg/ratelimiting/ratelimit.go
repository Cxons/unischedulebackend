package ratelimiting

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)



func Limit(maxRequests int, duration time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		maxRequests,
		duration,
		httprate.WithKeyFuncs(httprate.KeyByIP),
	)
}




