package domain

type ChallengeType int

const (
	ChallengeTypeHashcash ChallengeType = iota // 0
	ChallengeTypeMBound                        // 1 // TODO: implement memory bound challenge
)
