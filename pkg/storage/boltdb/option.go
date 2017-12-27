package boltdb

var (
	defaultBucket = []byte("_winston")
)

// Option db option
type Option func(s *Storage)

// Bucket modify bucket name
func Bucket(name []byte) Option {
	return func(s *Storage) {
		if len(name) > 0 {
			s.bucket = name
		}
	}
}
