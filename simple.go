package cbor

import (
	"bytes"
	"fmt"
)

// DecodeSimple decodes the second byte of a simple cbor value - no values are defined in this range in RFC7049 so this will always cause an error. This behavior is not strictly conformant to RFC7049 but there is no easy way to return the raw value without it being ambigous.
func DecodeSimple(fb byte, r *bytes.Reader) (interface{}, error) {
	fb = fb & 0x1F
	return nil, fmt.Errorf("simple value %d is not assigned in RFC7049", fb)
}
