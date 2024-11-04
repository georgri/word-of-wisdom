package hashcash

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/georgri/word-of-wisdom/internal/domain"
	"github.com/georgri/word-of-wisdom/pkg/util"
	"github.com/google/uuid"
	"math"
	"time"
)

const (
	defaultChallengeType = domain.ChallengeTypeHashcash
)

type service struct {
}

// SolveChallenge searches for the Solution that makes the hash of the challenge good
func (s *service) SolveChallenge(ctx context.Context, challenge *domain.HashcashChallenge) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if util.CheckLeadingZeroBits(challenge.LeadingZeroBits, sha256.Sum256([]byte(s.GetHashKey(challenge)))) {
				return
			}
			// likely to timeout sooner
			if challenge.Solution == math.MaxUint {
				return fmt.Errorf("failed to calculate the counter")
			}
			challenge.Solution++
		}
	}
}

// GetHashKey makes the challenge key to be hashed and checked
func (s *service) GetHashKey(challenge *domain.HashcashChallenge) string {
	return fmt.Sprintf(
		"%v:%v:%v:%v:%v:%v",
		challenge.ID,
		challenge.Type,
		challenge.LeadingZeroBits,
		challenge.UnixTime,
		challenge.ResourceID,
		challenge.Solution,
	)
}

// GenerateChallenge makes a new hashcash challenge
func (s *service) GenerateChallenge(leadingZeroBits uint, resourceID string) (*domain.HashcashChallenge, error) {
	return &domain.HashcashChallenge{
		Type:            defaultChallengeType,
		LeadingZeroBits: leadingZeroBits,
		UnixTime:        time.Now().Unix(),
		ResourceID:      resourceID,
		ID:              uuid.New().String(),
	}, nil
}

// NewService returns a new instance of the hashcash service
func NewService() *service {
	return &service{}
}
