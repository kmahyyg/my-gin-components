package mem_ratelimiter

import (
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	ramstore "github.com/ulule/limiter/v3/drivers/store/memory"
)

// LimiterBuild parse user input limit string and then return middleware of gin
// by default, it uses in-memory store
// @param lmt: by default, use 200 reqs/minute, "200-M", string
// @param bizName: business name, as prefix
// @return middleware: gin.HandlerFunc, trust X-Real-IP and X-Forwarded-For
func LimiterBuild(lmt string, bizName string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(lmt)
	if err != nil {
		panic(err)
	}
	store := ramstore.NewStoreWithOptions(
		limiter.StoreOptions{
			Prefix:          bizName + "-",
			CleanUpInterval: 3600,
		},
	)
	lmtInstance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	lmtMidWare := mgin.NewMiddleware(lmtInstance)
	return lmtMidWare
}
