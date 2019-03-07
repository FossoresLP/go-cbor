package cbor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// EncodeFloat encodes a 64-bit floating point value to a cbor value
func EncodeFloat(in float64) []byte {
	u := math.Float64bits(in)
	out := make([]byte, 9)
	out[0] = 0xFB
	binary.BigEndian.PutUint64(out[1:], u)
	return out
}

// DecodeFloat16 decodes a cbor value to a 64-bit floating point value (the cbor value was 16-bit though)
func DecodeFloat16(_ byte, r *bytes.Reader) (interface{}, error) { // nolint: interfacer
	b := make([]byte, 2)
	n, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if n < 2 {
		return nil, fmt.Errorf("expected 2 bytes of data, got only %d", n)
	}
	// The following conversion has been adapted from the CBOR specification in RFC 7049
	half := (int(b[0]) << 8) + int(b[1])
	exp := (half >> 10) & 0x1f
	mant := half & 0x3ff
	var val float64
	if exp == 0 {
		val = math.Ldexp(float64(mant), -24)
	} else if exp != 31 {
		val = math.Ldexp(float64(mant+1024), exp-25)
	} else if mant == 0 {
		val = math.Inf(0)
	} else {
		val = math.NaN()
	}
	if half&0x8000 != 0 {
		return -val, nil
	}
	return val, nil
}

// DecodeFloat32 decodes a cbor value to a 32-bit floating point value
func DecodeFloat32(_ byte, r *bytes.Reader) (interface{}, error) { // nolint: interfacer
	b := make([]byte, 4)
	n, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if n < 4 {
		return nil, fmt.Errorf("expected 4 bytes of data, got only %d", n)
	}
	ui := binary.BigEndian.Uint32(b)
	return math.Float32frombits(ui), nil
}

// DecodeFloat64 decodes a cbor value to a 64-bit floating point value
func DecodeFloat64(_ byte, r *bytes.Reader) (interface{}, error) { // nolint: interfacer
	b := make([]byte, 8)
	n, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if n < 8 {
		return nil, fmt.Errorf("expected 8 bytes of data, got only %d", n)
	}
	ui := binary.BigEndian.Uint64(b)
	return math.Float64frombits(ui), nil
}
