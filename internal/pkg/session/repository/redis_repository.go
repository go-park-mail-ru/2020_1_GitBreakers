package repository

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type SessionRedis struct {
	redisDB *redis.Conn
}

func (redisConn *SessionRedis) Create(id string, login string, expire time.Duration) error {
	//insert to redis
	return nil
}

func (redisConn *SessionRedis) Delete(id string) error {
	//delete from redis
	return nil
}

func (redisConn *SessionRedis) GetLogin(id string) error {
	//insert to redis
	return nil
}
