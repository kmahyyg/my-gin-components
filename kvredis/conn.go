package kvredis

import (
	"errors"
	"github.com/kmahyyg/my-gin-components/common-conf"
	"github.com/go-redis/redis/v8"
	"context"
)

type RedisConnectionFactory struct {
	isBuilt bool
	rClient *redis.Client
}

func (rcf *RedisConnectionFactory) GetRedisConn(conf common_conf.RedisConfig) (*redis.Client,error) {
	if rcf.isBuilt && rcf.rClient != nil {return rcf.rClient,nil}
	rcf.rClient = redis.NewClient(&redis.Options{
		Addr:       conf.Addr,
		Username:   conf.Username,
		Password:   conf.Password,
		DB:         conf.DBNum,
		MaxRetries: 5,
	})
	err := rcf.CheckRedisConn()
	if err != nil {return nil, err}
	rcf.isBuilt = true
	return rcf.rClient, nil
}

func (rcf *RedisConnectionFactory) CheckRedisConn() error {
	if !rcf.isBuilt || rcf.rClient == nil {return errors.New("no redis client is connected to server")}
	ctx := context.Background()
	_, err := rcf.rClient.Ping(ctx).Result()
	return err
}

func (rcf *RedisConnectionFactory) ResetRedisConf() {
	rcf.isBuilt = false
	rcf.rClient = nil
	return
}
