package domain

type HashcashChallenge struct {
	ID              string        `json:"id"`
	Type            ChallengeType `json:"type"`
	LeadingZeroBits uint          `json:"leading_zero_bits"`
	UnixTime        int64         `json:"unix_time"`
	ResourceID      string        `json:"resource_id"`
	Solution        uint64        `json:"solution"`
}
