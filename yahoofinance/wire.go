package yahoofinance

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

var (
	ErrNilValue        = fmt.Errorf("nil value")
	ErrFieldOutOfRange = fmt.Errorf("field out of range")
)

const (
	WireTypeFloat32 = 5
	WireTypeString  = 2
)

func unmarshalWireMessage(data []byte, v any) error {
	if v == nil {
		return ErrNilValue
	}

	valueOf := reflect.ValueOf(v).Elem()
	typeOf := valueOf.Type()

	protoIndexes := map[int]reflect.Value{}
	fieldCount := typeOf.NumField()

	for i := range fieldCount {
		field := typeOf.Field(i)
		protoIndex := field.Tag.Get("protoIndex")
		if protoIndex == "" {
			return fmt.Errorf("field %s has no protoIndex tag", field.Name)
		}
		n, err := strconv.Atoi(protoIndex)
		if err != nil {
			return fmt.Errorf("invalid protoIndex tag for field %s: %w", field.Name, err)
		}

		protoIndexes[n] = valueOf.Field(i)
	}

	i := 0
	for i < len(data) {
		if i >= len(data) {
			break
		}

		tag := data[i]
		fieldNum := int(tag >> 3)

		field, ok := protoIndexes[fieldNum]
		i++
		if !ok {
			continue
		}
		wireType := int(tag & 0x07)

		switch wireType {
		case WireTypeString:
			length, n := decodeVarint(data[i:])
			i += n
			if i+int(length) <= len(data) {
				if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
					strVal := string(data[i : i+int(length)])
					field.SetString(strVal)
				}
				i += int(length)
			}
		case WireTypeFloat32:
			if i+4 <= len(data) {
				val := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
				floatVal := *(*float32)(unsafe.Pointer(&val))

				if field.IsValid() && field.CanSet() && field.Kind() == reflect.Float32 {
					field.SetFloat(float64(floatVal))
				}
				i += 4
			}
		default:
			continue
		}
	}
	return nil
}

func decodeVarint(data []byte) (uint64, int) {
	var result uint64
	var shift uint
	var i int

	for i = 0; i < len(data) && i < 10; i++ {
		b := data[i]
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			return result, i + 1
		}
		shift += 7
	}
	return result, i
}
