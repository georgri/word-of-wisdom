package config

import "time"

type Config struct {
	ServerHost          string        `envconfig:"SERVER_HOST" default:":13371"`
	ChallengeComplexity uint          `envconfig:"CHALLENGE_COMPLEXITY" default:"12"`
	SolutionTimeout     time.Duration `envconfig:"SOLUTION_TIMEOUT" default:"10s"`
	ReadTimeout         time.Duration `envconfig:"READ_TIMEOUT" default:"30s"`
}
