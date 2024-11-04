package hashcash

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/georgri/word-of-wisdom/internal/domain"
	"github.com/georgri/word-of-wisdom/pkg/util"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type hashcashTestSuite struct {
	suite.Suite

	ctx               context.Context
	currentGoroutines goleak.Option

	// testing data
	challengeJSON []byte
	hashKey       string
	challenge     *domain.HashcashChallenge
	service       *service
}

func (s *hashcashTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.currentGoroutines = goleak.IgnoreCurrent()
	s.challengeJSON = []byte(`{"id":"53355a7b-dd66-48c9-82dd-6c50e6c1a990","type":0,"leading_zero_bits":12,"unix_time":1730685507,"resource_id":"0fcc98ab-dc7e-464c-8ef5-72fd725ad86c"}`)
	s.hashKey = "53355a7b-dd66-48c9-82dd-6c50e6c1a990:0:12:1730685507:0fcc98ab-dc7e-464c-8ef5-72fd725ad86c:0"
	s.service = NewService()
	s.challenge = &domain.HashcashChallenge{
		ID:              "53355a7b-dd66-48c9-82dd-6c50e6c1a990",
		Type:            domain.ChallengeTypeHashcash,
		LeadingZeroBits: 12,
		UnixTime:        1730685507,
		ResourceID:      "0fcc98ab-dc7e-464c-8ef5-72fd725ad86c",
	}
}

func (s *hashcashTestSuite) TearDownTest() {
	goleak.VerifyNone(s.T(), s.currentGoroutines)
}

func (s *hashcashTestSuite) TestHashKey() {
	hashKey := s.service.GetHashKey(s.challenge)
	s.EqualValues(s.hashKey, hashKey)
}

func (s *hashcashTestSuite) TestParseJSON() {
	parsedChallenge := &domain.HashcashChallenge{}
	err := json.Unmarshal(s.challengeJSON, parsedChallenge)
	s.Nil(err)
	s.Equal(s.challenge, parsedChallenge)

	hashKey := s.service.GetHashKey(s.challenge)
	s.EqualValues(s.hashKey, hashKey)
}

func (s *hashcashTestSuite) TestSolveChallenge() {
	err := s.service.SolveChallenge(s.ctx, s.challenge)
	s.Nil(err)

	hashKey := s.service.GetHashKey(s.challenge)
	hash := sha256.Sum256([]byte(hashKey))

	s.True(util.CheckLeadingZeroBits(s.challenge.LeadingZeroBits, hash))
	s.T().Logf("calculated hash result is %x\n", hash)
}

func TestHashcashTestSuite(t *testing.T) {
	suite.Run(t, &hashcashTestSuite{})
}
