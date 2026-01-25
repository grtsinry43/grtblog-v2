package federation

import "errors"

var (
	ErrUnsupportedSignatureAlgorithm = errors.New("unsupported signature algorithm")
	ErrMissingSignatureHeader        = errors.New("missing signature header")
	ErrInvalidDigest                 = errors.New("invalid digest header")
	ErrSignatureExpired              = errors.New("signature timestamp expired")
)
