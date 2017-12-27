package storage

import "bytes"

// PrefixEnd key upper bound for the prefix
func PrefixEnd(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}

	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			end := make([]byte, i+1)
			copy(end, prefix)
			end[i] = c + 1
			return end
		}
	}

	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0x00
	return end
}

// KeyInRange check if key in the range [start, end)
func KeyInRange(key, start, end []byte) bool {
	if start != nil && bytes.Compare(key, start) < 0 {
		return false
	}

	if end != nil && bytes.Compare(key, end) >= 0 {
		return false
	}

	return true
}
