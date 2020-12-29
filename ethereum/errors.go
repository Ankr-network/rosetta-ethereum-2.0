package ethereum

import "errors"

// Client errors
var (
	ErrBlockOrphaned         = errors.New("block orphaned")
	ErrCallParametersInvalid = errors.New("call parameters invalid")
	ErrCallOutputMarshal     = errors.New("call output marshal")
	ErrCallMethodInvalid     = errors.New("call method invalid")
	ErrBlockNotExists        = errors.New("block doesn't exist")
)
