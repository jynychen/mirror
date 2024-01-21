package mirror

import "errors"

var (
	ErrEmptySrcRepoURL = errors.New("source repository ur is required")
	ErrEmptyDstRepoURL = errors.New("destination repository url is required")
)
