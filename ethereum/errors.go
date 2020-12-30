package ethereum

import "errors"

// Client errors
var (
	ErrBlockOrphaned         = errors.New("block orphaned")
	ErrCallParametersInvalid = errors.New("call parameters invalid")
	ErrCallOutputMarshal     = errors.New("call output marshal")
	ErrCallMethodInvalid     = errors.New("call method invalid")
	ErrBlockNotFound         = errors.New("block not found")
	ErrBlockMissed           = errors.New("block is missed")
)
