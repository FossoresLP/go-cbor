package cbor

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
)

// Encoder is an interface for types that know how to turn themselfes into cbor
type Encoder interface {
	EncodeCBOR() []byte
}

// Marshaler is an interface (derived from encoding.*Marshaler) for types that know how to turn themselfes into cbor (encoding will panic when error is not nil)
type Marshaler interface {
	MarshalCBOR() ([]byte, error)
}

var (
	tagMap map[uint64]func(interface{}) (interface{}, error)

	// Types used to check for interface implementations
	cborEncoder     reflect.Type = reflect.TypeOf((*Encoder)(nil)).Elem()
	cborMarshaler   reflect.Type = reflect.TypeOf((*Marshaler)(nil)).Elem()
	binaryMarshaler reflect.Type = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	textMarshaler   reflect.Type = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func init() {
	tagMap = make(map[uint64]func(interface{}) (interface{}, error))
}

// Encode encodes an arbitrary type or data structure as a cbor value
func Encode(in interface{}) []byte {
	return encode(in)
}

func encode(in interface{}) []byte {
	val := reflect.ValueOf(in)
	t := reflect.TypeOf(in)
	if t.Implements(cborEncoder) {
		return in.(Encoder).EncodeCBOR()
	}
	if t.Implements(cborMarshaler) {
		b, err := in.(Marshaler).MarshalCBOR()
		if err != nil {
			panic(err)
		}
		return b
	}
	if t.Implements(binaryMarshaler) {
		b, err := in.(encoding.BinaryMarshaler).MarshalBinary()
		if err != nil {
			panic(err)
		}
		return EncodeBytes(b)
	}
	switch val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return EncodeUint(val.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return EncodeInt(val.Int())
	case reflect.Float32, reflect.Float64:
		return EncodeFloat(val.Float())
	case reflect.String:
		return EncodeString(val.String())
	case reflect.Slice:
		if reflect.TypeOf(in) == reflect.TypeOf([]byte{0x00}) {
			return EncodeBytes(val.Bytes())
		}
		return EncodeArray(in)
	case reflect.Array:
		if val.Len() > 0 && reflect.TypeOf(val.Index(0)) == reflect.TypeOf([]byte{0x00}) {
			s := val.Slice(0, val.Len())
			return EncodeBytes(s.Bytes())
		}
		return EncodeArray(in)
	case reflect.Map:
		return EncodeMap(in)
	case reflect.Bool:
		if val.Bool() {
			return []byte{0xF5}
		}
		return []byte{0xF4}
	case reflect.Ptr, reflect.Interface:
		return encode(val.Elem().Interface())
	case reflect.Struct:
		return EncodeStruct(in)
	}
	if t.Implements(textMarshaler) {
		s, err := in.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			panic(err)
		}
		return EncodeString(string(s))
	}
	panic(fmt.Sprintf("encode fell trough: cannot encode type %s\nif it is a custom type make sure it implements a supported encoding interface\notherwise consult the README and open an issue if this type is not listed as unsupported there\n", reflect.TypeOf(in)))
}

// Decode decodes a cbor value to an arbitrary type or data structure
func Decode(in []byte) (interface{}, error) {
	r := bytes.NewReader(in)
	for fb, err := r.ReadByte(); err == nil; {
		return decode(fb, r)
	}
	return nil, nil
}

func decode(fb byte, r *bytes.Reader) (interface{}, error) {
	major := fb >> 5
	switch major {
	case 0:
		return DecodeUint(fb, r)
	case 1:
		return DecodeInt(fb, r)
	case 2:
		return DecodeBytes(fb, r)
	case 3:
		return DecodeString(fb, r)
	case 4:
		return DecodeArray(fb, r)
	case 5:
		return DecodeMap(fb, r)
	case 6:
		return DecodeTag(fb, r)
	case 7:
		additional := fb & 0x1F
		switch additional {
		case 20:
			return false, nil
		case 21:
			return true, nil
		case 22:
			return nil, nil
		case 23:
			return nil, nil
		case 24:
			return DecodeSimple(fb, r)
		case 25:
			return DecodeFloat16(fb, r)
		case 26:
			return DecodeFloat32(fb, r)
		case 27:
			return DecodeFloat64(fb, r)
		case 31:
			return nil, fmt.Errorf("unexpected break outside an indefinite length element")
		default:
			return nil, fmt.Errorf("simple value %d is not assigned in RFC7049", additional)
		}
	}
	return nil, nil // This will never be reached
}
