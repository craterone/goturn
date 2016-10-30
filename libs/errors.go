package libs

import "errors"

//stun
var (
	ERROR_INVALID_REQUEST = errors.New("invalid STUN request")
	ERROR_RFC3489 = errors.New("no magic cookie , no supported RFC3489")
)

//config
var (
	ERROR_NO_SUITABLE_CONFIG = errors.New("no suitable config")
)