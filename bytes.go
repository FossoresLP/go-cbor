package cbor

import (
	"bytes"
	"fmt"
)

// EncodeBytes encodes a byte slice to a cbor value
func EncodeBytes(b []byte) []byte {
	u := EncodeUint(uint64(len(b)))
	u[0] = u[0] | 0x40 // Major type 2
	return append(u, b...)
}

// EncodeString encodes a string to a cbor value
func EncodeString(s string) []byte {
	u := EncodeUint(uint64(len(s)))
	u[0] = u[0] | 0x60 // Major type 3
	return append(u, []byte(s)...)
}

// DecodeBytes decodes a cbor value to a byte slice
func DecodeBytes(fb byte, r *bytes.Reader) ([]byte, error) {
	fb = fb & 0x1F // Remove major type (0x00011111 = 0x1F) from first byte
	if fb == 31 {
		return decodeIndefiniteBytes(r)
	}
	l, err := DecodeUint(fb, r)
	if err != nil {
		return nil, err
	}
	out := make([]byte, l)
	n, err := r.Read(out)
	if err != nil {
		return nil, fmt.Errorf("reading %d bytes failed: %s", l, err.Error())
	}
	if uint64(n) < l {
		return nil, fmt.Errorf("expected to read %d bytes but got only %d", l, n)
	}
	return out, nil
}

func decodeIndefiniteBytes(r *bytes.Reader) ([]byte, error) {
	var out bytes.Buffer
	fb, err := r.ReadByte()
	for ; err == nil && (fb>>5 == 2); fb, err = r.ReadByte() {
		v, err := DecodeBytes(fb, r)
		if err != nil {
			return nil, err
		}
		out.Write(v)
	}
	if err != nil {
		return nil, err
	}
	if fb != 0xFF {
		return nil, fmt.Errorf("indefinite length byte string contains unexpected type %d with additional value %d", fb>>5, fb&0x1F)
	}
	return out.Bytes(), nil
}

// DecodeString decodes a cbor value to a string
func DecodeString(fb byte, r *bytes.Reader) (string, error) {
	fb = fb & 0x1F // Remove major type (0x00011111 = 0x1F) from first byte
	if fb == 31 {
		return decodeIndefiniteString(r)
	}
	l, err := DecodeUint(fb, r)
	if err != nil {
		return "", err
	}
	out := make([]byte, l)
	n, err := r.Read(out)
	if err != nil {
		return "", fmt.Errorf("reading %d bytes for string failed: %s", l, err.Error())
	}
	if uint64(n) < l {
		return "", fmt.Errorf("string indicates a length of %d, but only %d bytes are left", l, n)
	}
	return string(out), nil
}

func decodeIndefiniteString(r *bytes.Reader) (string, error) {
	var out bytes.Buffer
	fb, err := r.ReadByte()
	for ; err == nil && (fb>>5 == 3); fb, err = r.ReadByte() {
		v, err := DecodeString(fb, r)
		if err != nil {
			return "", err
		}
		out.WriteString(v)
	}
	if err != nil {
		return "", err
	}
	if fb != 0xFF {
		return "", fmt.Errorf("indefinite length string contains unexpected type %d with additional value %d", fb>>5, fb&0x1F)
	}
	return out.String(), nil
}
