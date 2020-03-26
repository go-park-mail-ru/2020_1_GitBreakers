package repository

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"time"
)

type SessionRedis struct {
	RedisDB *redis.Conn
}

func (repo *SessionRedis) Create(sid string, login string, expire time.Duration) (string, error) {
	if err := repo.RedisDB.Set(sid, login, expire).Err(); err != nil {
		return "", errors.Wrap(err, "session_repository: redis Set failed for")
	}
	return sid, nil
}

func (repo *SessionRedis) GetLoginById(sessionId string) (string, error) {
	login, err := repo.RedisDB.Get(sessionId).Result()

	switch {
	case err == redis.Nil:
		return "", errors.WithMessagef(entityerrors.DoesNotExist(),
			"session_repository: key does not exist, sessionId=%v", sessionId)
	case err != nil:
		return "", errors.Wrapf(err,
			"session_repository: redis Get failed with sessionId=%v", sessionId)
	}

	return login, nil
}

func (repo *SessionRedis) DeleteById(sid string) error {
	deletedKeysCount, err := repo.RedisDB.Del(sid).Result()
	switch {
	case err != nil:
		return errors.Wrapf(err, "session_repository: redis Del failed with sid=%v", sid)
	case deletedKeysCount != 1:
		return errors.WithMessagef(entityerrors.DoesNotExist(),
			"session_repository: key does not exist, sid=%v", sid)
	}
	return nil
}
