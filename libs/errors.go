package libs

import "errors"

//stun
var (
	ERROR_INVALID_REQUEST = errors.New("Invalid STUN request")
	ERROR_RFC3489 = errors.New("no magic cookie , no supported RFC3489")
)

//config
var (
	ERROR_NO_SUITABLE_CONFIG = errors.New("No suitable config")
)