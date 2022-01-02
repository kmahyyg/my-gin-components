package ssbam

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kmahyyg/my-gin-components/common-conf"
	"github.com/kmahyyg/my-gin-components/kvredis"
	"net/http"
	"time"
)

const (
	cookie_SESSION_NAME     = "SESSIONID"
	gincontext_SESSION_USER = "sess_u"
	redis_SESSION_PREFIX    = "sess_"
)

func SessionAuthentication(rcf *kvredis.RedisConnectionFactory, conf common_conf.WebConfig) gin.HandlerFunc {
	rCli, err := rcf.GetRedisConn()
	if err != nil {
		panic(err)
	}
	return func(gctx *gin.Context) {
		// global var
		ctx := context.Background()
		resetFlag := false
		// cookie in request
		ck, err := gctx.Request.Cookie(cookie_SESSION_NAME)
		// if cookie is not nil, check if exists in redis
		// including: - session id exist in cookie but invalidated
		//            - session id exist and valid
		if ck != nil {
			_, err = uuid.Parse(ck.Value)
			if err != nil {
				// session id is not a valid uuidv4
				resetFlag = true
			}
			if !resetFlag {
				// check if session is logged in
				val, err := rCli.Get(ctx, redis_SESSION_PREFIX+ck.Value).Result()
				if err != nil {
					// not logged in, reset
					resetFlag = true
				} else {
					// logged in, set context for convenience as authentication.
					// note: this system is allowed to login at multiple location.
					gctx.Set(cookie_SESSION_NAME, ck.Value)
					gctx.Set(gincontext_SESSION_USER, val) // only exists if bind to a user
				}
			}
		}
		// new request from unknown user or invalid user
		if err == http.ErrNoCookie || err != nil || resetFlag {
			sessionNew := uuid.New().String()
			gctx.Set(cookie_SESSION_NAME, sessionNew)
			// current session is not binded to any user
			gctx.SetCookie(cookie_SESSION_NAME, sessionNew, 3600, "/", conf.DomainName, false, true)
			err = rCli.Set(ctx, redis_SESSION_PREFIX+sessionNew, "", time.Second*3600).Err()
			if err != nil {
				panic(err)
			}
		}
		gctx.Next()
	}
}
