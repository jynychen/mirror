package mirror

import "errors"

var (
	ErrEmptySourceRepoURL  = errors.New("source repository ur is required")
	ErrEmptyDestinationURL = errors.New("destination repository url is required")
)
