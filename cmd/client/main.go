package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/georgri/word-of-wisdom/internal/domain"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/rs/zerolog"

	"github.com/georgri/word-of-wisdom/internal/config"
	hashcashservice "github.com/georgri/word-of-wisdom/pkg/hashcash"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	var (
		cfg    config.Config
		logger = log.New(os.Stdout).With().Timestamp().Logger()
	)
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Err(err).Msg("failed to parse config")
		return
	}
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", cfg.ServerHost)
	if err != nil {
		logger.Err(err).Msg("failed to connect to the server")
		return
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			logger.Err(err).Msg("failed to close the connection")
		}
	}(conn)

	// init challenge
	_, err = conn.Write([]byte("{}"))
	if err != nil {
		logger.Err(err).Msg("failed to init challenge")
		return
	}

	// get the challenge data
	reader := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	if err != nil {
		logger.Err(err).Msg("failed to set the read timeout")
	}
	challenge := &domain.HashcashChallenge{}
	if err = json.NewDecoder(reader).Decode(challenge); err != nil {
		logger.Err(err).Msg("failed to get the challenge data from the server")
		return
	}

	// compute the solution and send the result
	svc := hashcashservice.NewService()
	ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, cfg.SolutionTimeout)
	defer cancelFunc()
	err = svc.SolveChallenge(ctxWithTimeout, challenge)
	if err != nil {
		logger.Err(err).Msg("failed to compute a solution in time")
		return
	}

	// send a solution
	conn, err = dialer.DialContext(ctx, "tcp", cfg.ServerHost)
	if err != nil {
		logger.Err(err).Msg("failed to connect to the server")
		return
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			logger.Err(err).Msg("failed to close the connection")
		}
	}(conn)

	writer := bufio.NewWriter(conn)
	err = json.NewEncoder(writer).Encode(challenge)
	if err != nil {
		logger.Err(err).Msg("failed to send a solution to the challenge ")
		return
	}

	// read the response
	reader = bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	if err != nil {
		logger.Warn().Err(err).Msg("failed to set read timeout")
	}
	var result string
	result, err = reader.ReadString('\n')
	if err != nil {
		logger.Err(err).Msg("failed to read a quote")
		return
	}
	logger.Info().Msgf("received a quote: %v", strings.TrimSpace(result))
}
