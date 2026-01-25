package federation

import "errors"

var (
	ErrFederationConfigNotFound   = errors.New("federation config not found")
	ErrFederationInstanceNotFound = errors.New("federation instance not found")
)
