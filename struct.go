package cbor

import (
	"bytes"
	"reflect"
	"strings"
)

// EncodeStruct encodes a struct (treated as a map) to a cbor value
func EncodeStruct(in interface{}) []byte {
	v := reflect.ValueOf(in)
	t := reflect.TypeOf(in)
	length := v.NumField()

	out := bytes.NewBuffer([]byte{0xBF}) // Major type 5 additional 31 - Indefinite length map - used due to the length including unexported fields
	for i := 0; i < length; i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		f := v.Field(i).Interface()
		s := t.Field(i)
		if tag := s.Tag.Get("cbor"); tag != "" {
			tags := strings.Split(tag, ",")
			out.Write(EncodeString(tags[0]))
		} else {
			out.Write(EncodeString(s.Name))
		}
		out.Write(encode(f))
	}
	out.WriteByte(0xFF)
	return out.Bytes()
}
