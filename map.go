package cbor

import (
	"bytes"
	"fmt"
	"reflect"
)

// EncodeMap encodes a map using arbitrary types as keys and values to a cbor value
func EncodeMap(in interface{}) []byte {
	val := reflect.ValueOf(in)
	length := val.Len()
	l := EncodeUint(uint64(length))
	l[0] = l[0] | 0xA0 // Major type 5
	out := bytes.NewBuffer(l)
	iter := val.MapRange()
	for iter.Next() {
		out.Write(encode(iter.Key().Interface()))
		out.Write(encode(iter.Value().Interface()))
	}
	return out.Bytes()
}

// DecodeMap decodes a cbor value to a map using arbitrary types as keys and values
func DecodeMap(fb byte, r *bytes.Reader) (map[interface{}]interface{}, error) {
	fb = fb & 0x1F
	if fb == 31 {
		return decodeIndefiniteMap(r)
	}
	l, err := DecodeUint(fb, r)
	if err != nil {
		return nil, err
	}
	out := make(map[interface{}]interface{}, l)
	for i := 0; i < int(l); i++ {
		fb, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("could not read key: %s", err.Error())
		}
		key, err := decode(fb, r)
		if err != nil {
			return nil, err
		}
		fb, err = r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("could not read value: %s", err.Error())
		}
		value, err := decode(fb, r)
		if err != nil {
			return nil, err
		}
		out[key] = value
	}
	return out, nil
}

func decodeIndefiniteMap(r *bytes.Reader) (map[interface{}]interface{}, error) {
	out := make(map[interface{}]interface{})
	fb, err := r.ReadByte()
	for ; err == nil && fb != 0xFF; fb, err = r.ReadByte() {
		key, err := decode(fb, r)
		if err != nil {
			return nil, err
		}
		fb, err = r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("could not read value: %s", err.Error())
		}
		value, err := decode(fb, r)
		if err != nil {
			return nil, err
		}
		out[key] = value
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse indefinite length map: %s", err.Error())
	}
	return out, nil
}
