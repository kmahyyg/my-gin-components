package redis_ratelimiter

import (
	"github.com/gin-gonic/gin"
	"github.com/kmahyyg/my-gin-components/kvredis"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// LimiterBuild parse user input limit string and then return middleware of gin
// by default, it uses redis store for distributed system.
// @param lmt: by default, use 200 reqs/minute, "200-M", string
// @param bizName: business name, as prefix
// @param rcf: RedisConnectionFactory, to build redis storage backend client
// @return middleware: gin.Handlerfunc, trust X-Real-IP and XFF
func LimiterBuild(lmt string, bizName string, rcf *kvredis.RedisConnectionFactory) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(lmt)
	if err != nil {panic(err)}
	redisBackend, err := rcf.GetRedisConn()
	if err != nil || redisBackend == nil {panic(err)}
	store, err := redisstore.NewStoreWithOptions(redisBackend, limiter.StoreOptions{
		Prefix:          bizName,
		CleanUpInterval: 3600,
	})
	lmtInstance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	lmtMidWare := mgin.NewMiddleware(lmtInstance)
	return lmtMidWare
}
