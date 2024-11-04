package randomquote

import (
	"math/rand"
	"time"
)

type service struct {
	quotes    []string
	generator *rand.Rand
}

// GetQuote returns a random quote from the quotes list
func (s *service) GetQuote() string {
	return s.quotes[s.generator.Intn(len(s.quotes))]
}

// NewService initializes a new instance
func NewService(quotes []string) *service {
	return &service{
		quotes:    quotes,
		generator: rand.New(rand.NewSource(time.Now().Unix())),
	}
}
