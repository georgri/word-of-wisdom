package util

// CheckLeadingZeroBits check hash contains at least leadingZeroBits amount of leading zero bits
func CheckLeadingZeroBits(leadingZeroBits uint, hash [32]byte) bool {
	if leadingZeroBits > 256 {
		leadingZeroBits = 256
	}
	for i := 0; i < int(leadingZeroBits); i++ {
		if hash[i>>3]&(1<<(7-i&7)) != 0 {
			return false
		}
	}
	return true
}
