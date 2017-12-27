package meta

import (
	"net"

	"go.uber.org/zap"
)

const (
	// MuxHeader tcp mux header
	MuxHeader byte = 1
)

// Service meta service
type Service struct {
	ln     net.Listener
	logger *zap.SugaredLogger
}

// Start start the meta service
func (s *Service) Start() error {
	return nil
}

// Close close the meta service
func (s *Service) Close() error {
	return nil
}

// WithLogger setup logger
func (s *Service) WithLogger(logger *zap.Logger) {
	if logger != nil {
		s.logger = logger.Sugar()
	}
}
