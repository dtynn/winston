package codes

import (
	"errors"

	"github.com/dtynn/winston/internal/pb"
)

// Err for data node management
var (
	ErrRaftDuplicateGroupID = &codeError{
		code:  pb.ResultCode_ResultCodeRaftDuplicateGroupID,
		cause: errors.New("raft: duplicate group id"),
	}
)
