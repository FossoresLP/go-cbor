package cbor

import (
	"bytes"
	"fmt"
	"reflect"
)

// EncodeArray encodes an array of arbitrary types to a cbor value
func EncodeArray(arr interface{}) []byte {
	val := reflect.ValueOf(arr)
	length := val.Len()
	l := EncodeUint(uint64(length))
	l[0] = l[0] | 0x80 // Major type 4
	out := bytes.NewBuffer(l)
	for i := 0; i < length; i++ {
		out.Write(encode(val.Index(i)))
	}
	return out.Bytes()
}

// DecodeArray decodes a cbor value to an array of arbitrary types
func DecodeArray(fb byte, r *bytes.Reader) ([]interface{}, error) {
	fb = fb & 0x1F
	if fb == 31 {
		return decodeIndefiniteArray(r)
	}
	l, err := DecodeUint(fb, r)
	if err != nil {
		return nil, err
	}
	out := make([]interface{}, l)
	for i := 0; i < int(l); i++ {
		fb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("expected %d more element(s) in array but input ends after item %d", int(l)-i-1, i+1)
		}
		v, err := decode(fb, r)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func decodeIndefiniteArray(r *bytes.Reader) ([]interface{}, error) {
	out := make([]interface{}, 0, 10)
	fb, err := r.ReadByte()
	for ; err == nil && fb != 0xFF; fb, err = r.ReadByte() {
		v, err := decode(fb, r)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse indefinite length array: %s", err.Error())
	}
	return out, nil
}
