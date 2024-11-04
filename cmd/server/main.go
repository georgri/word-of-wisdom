package main

import (
	"context"
	"github.com/georgri/word-of-wisdom/pkg/quotes"
	log "github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"

	"github.com/georgri/word-of-wisdom/internal/config"
	serverpkg "github.com/georgri/word-of-wisdom/internal/server"
	"github.com/georgri/word-of-wisdom/internal/storage"
	hashcashservice "github.com/georgri/word-of-wisdom/pkg/hashcash"
	"github.com/georgri/word-of-wisdom/pkg/randomquote"
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

	quoteList := quotes.GetHardcodedQuotes()
	server := serverpkg.NewServer(
		uuid.New().String(),
		cfg.ServerHost,
		cfg.ChallengeComplexity,
		cfg.SolutionTimeout,
		cfg.ReadTimeout,
		storage.NewStorage(),
		randomquote.NewService(quoteList),
		hashcashservice.NewService(),
		logger,
	)

	// run the server
	listenErrors := make(chan error, 1)
	go func() {
		err := server.Serve(ctx)
		if err != nil {
			listenErrors <- err
		}
	}()
	logger.Info().Msg("server started")
	select {
	case <-ctx.Done():
	case err := <-listenErrors:
		logger.Err(err).Msg("server failed")
		return
	}
	logger.Info().Msg("server stopped")
}
