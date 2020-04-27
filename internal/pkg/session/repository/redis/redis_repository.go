package redis

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-redis/redis/v7"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type SessionRedis struct {
	redisDb    *redis.Client
	basePrefix string
}

func NewSessionRedis(client *redis.Client, prefix string) SessionRedis {
	return SessionRedis{client, prefix}
}

func (repo *SessionRedis) createKey(sessionId string) string {
	return repo.basePrefix + sessionId
}

func (repo *SessionRedis) convertToString(session models.Session) (string, error) {
	sessionJSON, err := easyjson.Marshal(session)
	if err != nil {
		return "", errors.Wrapf(err, "session_repository: JSON Marshal failed for %+v", session)
	}
	return string(sessionJSON), nil
}

func (repo *SessionRedis) convertFromString(sessionString string) (models.Session, error) {
	var storedSession models.Session
	if err := easyjson.Unmarshal([]byte(sessionString), &storedSession); err != nil {
		return storedSession, errors.Wrapf(err,
			"session_repository: JSON Unmarshal failed with sessionString=%v", sessionString)
	}
	return storedSession, nil
}

func (repo *SessionRedis) Create(session models.Session, expire time.Duration) (string, error) {
	session.ID = uuid.NewV4().String()
	key := repo.createKey(session.ID)
	sessionJSON, err := repo.convertToString(session)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if err := repo.redisDb.Set(key, sessionJSON, expire).Err(); err != nil {
		return "", errors.Wrapf(err, "session_repository: redis Set failed for %+v", session)
	}
	return session.ID, nil
}

func (repo *SessionRedis) GetSessByID(sessionId string) (models.Session, error) {
	storedSession := models.Session{}
	key := repo.createKey(sessionId)
	sessionJSON, err := repo.redisDb.Get(key).Result()
	switch {
	case err == redis.Nil:
		return storedSession, errors.WithMessagef(entityerrors.DoesNotExist(),
			"session_repository: key does not exist, sessionId=%v", sessionId)
	case err != nil:
		return storedSession, errors.Wrapf(err,
			"session_repository: redis Get failed with sessionId=%v", sessionId)
	}

	if storedSession, err = repo.convertFromString(sessionJSON); err != nil {
		return storedSession, errors.WithStack(err)
	}

	return storedSession, nil
}

func (repo *SessionRedis) DeleteByID(sid string) error {
	deletedKeysCount, err := repo.redisDb.Del(repo.createKey(sid)).Result()
	switch {
	case err != nil:
		return errors.Wrapf(err, "session_repository: redis Del failed with sid=%v", sid)
	case deletedKeysCount != 1:
		return errors.WithMessagef(entityerrors.DoesNotExist(),
			"session_repository: key does not exist, sid=%v", sid)
	}
	return nil
}
