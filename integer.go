package cbor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// EncodeInt encodes an int64 as a cbor value
func EncodeInt(i int64) []byte {
	if i >= 0 {
		return EncodeUint(uint64(i))
	}
	i++
	out := EncodeUint(uint64(-i))
	out[0] = out[0] | 0x20 // Major type 1 0b001xxxxx = 0x20 = 32
	return out
}

// EncodeUint encodes an uint64 as a cbor value
func EncodeUint(i uint64) []byte {
	if i < 24 {
		return []byte{byte(i)}
	} else if i <= math.MaxUint8 {
		return []byte{24, byte(i)}
	} else if i <= math.MaxUint16 {
		out := make([]byte, 3)
		out[0] = 25
		binary.BigEndian.PutUint16(out[1:], uint16(i))
		return out
	} else if i <= math.MaxUint32 {
		out := make([]byte, 5)
		out[0] = 26
		binary.BigEndian.PutUint32(out[1:], uint32(i))
		return out
	}
	out := make([]byte, 9)
	out[0] = 27
	binary.BigEndian.PutUint64(out[1:], i)
	return out
}

// DecodeInt decodes a cbor value to an int64
func DecodeInt(fb byte, r *bytes.Reader) (int64, error) {
	fb = fb & 0x1F // Get first byte and remove major type (0x00011111 = 0x1F)
	u, err := DecodeUint(fb, r)
	if err != nil {
		return 0, err
	}
	if u > math.MaxInt64 {
		return 0, fmt.Errorf("-1 - %d does exceed the range of an int64", u)
	}
	return -1 - (int64(u)), nil
}

// DecodeUint decodes a cbor value to an uint64
func DecodeUint(fb byte, r *bytes.Reader) (uint64, error) {
	switch fb {
	case 24:
		b, err := r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("error reading additional bytes: %s", err.Error())
		}
		return uint64(b), nil
	case 25:
		b := make([]byte, 2)
		n, err := r.Read(b)
		if err != nil {
			return 0, fmt.Errorf("error reading additional bytes: %s", err.Error())
		}
		if n < 2 {
			return 0, fmt.Errorf("type indicates two additional bytes but only %d read", n)
		}
		return uint64(binary.BigEndian.Uint16(b)), nil
	case 26:
		b := make([]byte, 4)
		n, err := r.Read(b)
		if err != nil {
			return 0, fmt.Errorf("error reading additional bytes: %s", err.Error())
		}
		if n < 4 {
			return 0, fmt.Errorf("type indicates four additional bytes but only %d read", n)
		}
		return uint64(binary.BigEndian.Uint32(b)), nil
	case 27:
		b := make([]byte, 8)
		n, err := r.Read(b)
		if err != nil {
			return 0, fmt.Errorf("error reading additional bytes: %s", err.Error())
		}
		if n < 8 {
			return 0, fmt.Errorf("type indicates eight additional bytes but only %d read", n)
		}
		return binary.BigEndian.Uint64(b), nil
	case 28, 29, 30, 31:
		return 0, fmt.Errorf("%d is not a valid value for additional information", fb)
	default:
		return uint64(fb), nil
	}
}
