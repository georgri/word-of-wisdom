package server

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	log "github.com/rs/zerolog"

	"github.com/georgri/word-of-wisdom/internal/domain"
	"github.com/georgri/word-of-wisdom/pkg/util"
)

type randomQuoteService interface {
	GetQuote() string
}

type hashcashService interface {
	GenerateChallenge(bits uint, resource string) (*domain.HashcashChallenge, error)
	SolveChallenge(ctx context.Context, data *domain.HashcashChallenge) (err error)
	GetHashKey(challenge *domain.HashcashChallenge) string
}

type challengeStorage interface {
	GetChallenge(id string) (challenge *domain.HashcashChallenge, ok bool)
	SetChallenge(id string, challenge *domain.HashcashChallenge)
	DeleteChallenge(id string)
}

type server struct {
	resourceID          string
	host                string
	randomQuoteService  randomQuoteService
	hashcashService     hashcashService
	challengeStorage    challengeStorage
	challengeComplexity uint
	solutionTimeout     time.Duration
	readTimeout         time.Duration
	logger              log.Logger
}

func NewServer(
	id string,
	addr string,
	challengeComplexity uint,
	solutionTimeout time.Duration,
	readTimeout time.Duration,
	storage challengeStorage,
	wisdomWordsService randomQuoteService,
	hashcashService hashcashService,
	logger log.Logger,
) *server {
	return &server{
		resourceID:          id,
		host:                addr,
		challengeComplexity: challengeComplexity,
		solutionTimeout:     solutionTimeout,
		readTimeout:         readTimeout,
		challengeStorage:    storage,
		randomQuoteService:  wisdomWordsService,
		hashcashService:     hashcashService,
		logger:              logger,
	}
}

// Serve contains the server logic and serves requests
func (s *server) Serve(ctx context.Context) (err error) {
	var listener net.Listener
	listener, err = net.Listen("tcp", s.host)
	if err != nil {
		return
	}
	defer func() {
		errClose := listener.Close()
		if errClose != nil {
			s.logger.Err(errClose).Msg("failed to close the listener")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var conn net.Conn
			conn, err = listener.Accept()
			if err != nil {
				s.logger.Err(err).Msg("failed to accept a connection")
				continue
			}
			go s.processConn(conn)
		}
	}
}

func (s *server) processConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			s.logger.Err(err).Msg("failed to close the connection")
		}
	}()
	reader := bufio.NewReader(conn)

	msgLine, err := reader.ReadString('\n')
	if err != nil {
		s.logger.Err(err).Msg("failed to read the request line")
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to set the read timeout")
	}

	// if the request is empty => it's a challenge request
	if len(strings.TrimSpace(msgLine)) == 0 {
		// create, save and return a new challenge
		var newChallenge *domain.HashcashChallenge
		newChallenge, err = s.hashcashService.GenerateChallenge(s.challengeComplexity, s.resourceID)
		if err != nil {
			s.logger.Err(err).Msg("failed to generate a hashcash challenge")
			return
		}
		s.challengeStorage.SetChallenge(newChallenge.ID, newChallenge)
		jsonChallenge, err := json.Marshal(newChallenge)
		if err != nil {
			s.logger.Err(err).Msg("failed to encode a challenge to JSON")
		}
		_, err = conn.Write(jsonChallenge)
		if err != nil {
			s.logger.Err(err).Msg("failed to return the newChallenge")
			return
		}
		return
	}

	challenge := &domain.HashcashChallenge{}
	err = json.Unmarshal([]byte(msgLine), challenge)
	if err != nil {
		s.logger.Err(err).Msg("failed to decode a challenge from the message")
		return
	}

	origChallenge, ok := s.challengeStorage.GetChallenge(challenge.ID)
	if !ok {
		s.logger.Err(err).Msg("orig challenge not found")
		return
	}
	defer s.challengeStorage.DeleteChallenge(challenge.ID)
	err = s.verifySolution(origChallenge, challenge)
	if err != nil {
		s.logger.Err(err).Msg("the challenge failed")
		return
	}
	// return the quote
	_, err = io.WriteString(conn, s.randomQuoteService.GetQuote()+"\n")
	if err != nil {
		s.logger.Err(err).Msg("failed to return the word of wisdom")
		return
	}
}

// verifySolution checks the challenge solution for validity and correctness
func (s *server) verifySolution(challenge, result *domain.HashcashChallenge) (err error) {
	// check that the challenge and the result are identical ignoring the solution part
	solution := result.Solution
	result.Solution = 0
	challenge.Solution = 0
	if s.hashcashService.GetHashKey(result) != s.hashcashService.GetHashKey(challenge) {
		return fmt.Errorf("invalid challenge result")
	}
	if time.Now().Unix()-challenge.UnixTime > int64(s.solutionTimeout.Seconds()) {
		return fmt.Errorf("challenge timeout expired")
	}
	// restore the solution and check the result
	result.Solution = solution
	hash := sha256.Sum256([]byte(s.hashcashService.GetHashKey(result)))
	if !util.CheckLeadingZeroBits(challenge.LeadingZeroBits, hash) {
		return fmt.Errorf("the challenge solution is incorrect")
	}
	s.logger.Info().Msgf("the challenge was solved with hash: %x", hash)
	return nil
}
