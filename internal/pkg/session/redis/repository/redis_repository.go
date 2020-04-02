package repository

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type SessionRedis struct {
	redisDb    *redis.Conn
	BasePrefix string
}

func NewSessionRedis(conn *redis.Conn, prefix string) SessionRedis {
	return SessionRedis{conn, prefix}
}

func (repo *SessionRedis) Create(session models.Session, expire time.Duration) (string, error) {
	session.Id = uuid.NewV4().String()
	key := repo.BasePrefix + session.Id
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", errors.Wrapf(err, "session_repository: JSON Marshal failed for %+v", session)
	}
	if err := repo.redisDb.Set(key, string(sessionJSON), expire).Err(); err != nil {
		return "", errors.Wrapf(err, "session_repository: redis Set failed for %+v", session)
	}
	return session.Id, nil
}

func (repo *SessionRedis) GetSessById(sessionId string) (models.Session, error) {
	storedSession := models.Session{}
	key := repo.BasePrefix + sessionId
	sessionJSON, err := repo.redisDb.Get(key).Result()
	switch {
	case err == redis.Nil:
		return storedSession, errors.WithMessagef(entityerrors.DoesNotExist(),
			"session_repository: key does not exist, sessionId=%v", sessionId)
	case err != nil:
		return storedSession, errors.Wrapf(err,
			"session_repository: redis Get failed with sessionId=%v", sessionId)
	}

	if err := json.Unmarshal([]byte(sessionJSON), &storedSession); err != nil {
		return storedSession, errors.Wrapf(err,
			"session_repository: JSON Unmarshal failed with sessionId=%v", sessionId)
	}

	return storedSession, nil
}

func (repo *SessionRedis) DeleteById(sid string) error {
	deletedKeysCount, err := repo.redisDb.Del(sid).Result()
	switch {
	case err != nil:
		return errors.Wrapf(err, "session_repository: redis Del failed with sid=%v", sid)
	case deletedKeysCount != 1:
		return errors.WithMessagef(entityerrors.DoesNotExist(),
			"session_repository: key does not exist, sid=%v", sid)
	}
	return nil
}
