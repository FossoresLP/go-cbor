package cbor

import (
	"bytes"
	"fmt"
)

// DecodeTag decodes a tag and following data element, if possible using a registered tag decoder
func DecodeTag(fb byte, r *bytes.Reader) (interface{}, error) {
	fb = fb & 0x1F
	val, err := DecodeUint(fb, r)
	if err != nil {
		return nil, err
	}
	nb, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("standalone tag is not allowed but failed to read following data: %s", err.Error())
	}
	data, err := decode(nb, r)
	if err != nil {
		return nil, err
	}
	if fn, ok := tagMap[val]; ok {
		return fn(data)
	}
	return data, nil
}
