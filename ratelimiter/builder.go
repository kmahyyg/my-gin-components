package ratelimiter

import (
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	ramlimitstore "github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiterBuild parse user input limit string and then return middleware of gin
// by default, it uses in-memory store
// @param lmt: by default, use 200 reqs/minute, "200-M", string
// @return middleware: gin.HandlerFunc
func RateLimiterBuild(lmt string, bizName string) gin.HandlerFunc {
	//
	rate, err := limiter.NewRateFromFormatted(lmt)
	if err != nil {
		panic(err)
	}
	store := ramlimitstore.NewStoreWithOptions(
		limiter.StoreOptions{
			Prefix:          bizName + "-",
			CleanUpInterval: 3600,
		},
	)
	lmtInstance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	lmtMidWare := mgin.NewMiddleware(lmtInstance)
	return lmtMidWare
}
