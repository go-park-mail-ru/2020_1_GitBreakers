package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

type redisTestSuite struct {
	suite.Suite
	redisRepository SessionRedis
	client          *redis.Client
	session         models.Session
}

func TestInit(t *testing.T) {
	suite.Run(t, new(redisTestSuite))
}

func (s *redisTestSuite) SetupSuite() {
	const redisTestPrefix = "test/"
	mr, err := miniredis.Run()
	if err != nil {
		s.Failf("setup redis test suite failed: ", " %+v\n", err)
	}
	s.client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	s.redisRepository = SessionRedis{basePrefix: redisTestPrefix, redisDb: s.client}
}

func (s *redisTestSuite) SetupTest() {
	s.session = models.Session{
		ID:     uuid.NewV4().String(),
		UserID: rand.Int(),
	}
}

func (s *redisTestSuite) TestCreate() {
	sessionId, err := s.redisRepository.Create(s.session, time.Minute)
	require.NoError(s.T(), err)

	sessStr := s.client.Get(s.redisRepository.createKey(sessionId)).Val()
	retSession, err := s.redisRepository.convertFromString(sessStr)
	require.NoError(s.T(), err)

	s.session.ID = retSession.ID

	require.Equal(s.T(), s.session, retSession)
}

func (s *redisTestSuite) TestGetSessById() {
	value, err := s.redisRepository.convertToString(s.session)
	require.NoError(s.T(), err)

	s.client.Set(s.redisRepository.createKey(s.session.ID), value, time.Hour)

	storedSession, err := s.redisRepository.GetSessByID(s.session.ID)
	require.NoError(s.T(), err)

	require.Equal(s.T(), s.session, storedSession)
}

func (s *redisTestSuite) TestGetSessByIdNegative() {
	_, err := s.redisRepository.GetSessByID(s.session.ID)
	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}

func (s *redisTestSuite) TestDeleteById() {
	value, err := s.redisRepository.convertToString(s.session)
	require.NoError(s.T(), err)

	s.client.Set(s.redisRepository.createKey(s.session.ID), value, time.Hour)
	require.NoError(s.T(), err)

	err = s.redisRepository.DeleteByID(s.session.ID)
	require.NoError(s.T(), err)
}

func (s *redisTestSuite) TestDeleteByIdNegative() {
	err := s.redisRepository.DeleteByID(s.session.ID)
	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}
