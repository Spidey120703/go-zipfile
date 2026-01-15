package serial

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func sizeof(v any) uint32 {
	if v == nil {
		return 0
	}

	if sizeofable, ok := v.(ISizeOf); ok {
		return sizeofable.SizeOf()
	}

	varType := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	switch varType.Kind() {
	case reflect.Struct:
		size := uint32(0)
		for i := 0; i < varType.NumField(); i++ {
			size += sizeof(value.Field(i).Interface())
		}
		return size
	case reflect.Pointer:
		if value.IsNil() {
			return 0
		}
		return sizeof(value.Elem())
	case reflect.Array, reflect.Slice:
		size := uint32(0)
		for j := 0; j < value.Len(); j++ {
			size += sizeof(value.Index(j).Interface())
		}
		return size
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		switch v.(type) {
		case bool, int8, uint8, *bool, *int8, *uint8:
			return 1
		case int16, uint16, *int16, *uint16:
			return 2
		case int32, uint32, *int32, *uint32:
			return 4
		case int64, uint64, *int64, *uint64:
			return 8
		case float32, *float32:
			return 4
		case float64, *float64:
			return 8
		}
	default:
		return 0
	}

	return 0
}

type marshaller struct {
	w io.WriteSeeker
}

func (m *marshaller) marshal(v any) (err error) {
	if v == nil {
		return
	}

	if serializable, ok := v.(ISerializable); ok {
		return serializable.Marshal(m.w)
	}

	varType := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	switch varType.Kind() {
	case reflect.Struct:
		for i := range varType.NumField() {
			if err = m.marshal(value.Field(i).Interface()); err != nil {
				return
			}
		}
	case reflect.Pointer:
		if value.IsNil() {
			return
		}
		return m.marshal(value.Elem())
	case reflect.Slice:
		if varType.Elem().Kind() == reflect.Uint8 && value.CanAddr() {
			// if value is bytes buffer, write it into file directly
			_, err = m.w.Write(value.Bytes())
			return
		}
		fallthrough
	case reflect.Array:
		for i := range value.Len() {
			if err = m.marshal(value.Index(i).Interface()); err != nil {
				return
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return binary.Write(m.w, binary.LittleEndian, value.Interface())
	default:
		return fmt.Errorf("unsupported type %v", varType)
	}

	return
}

func Marshal(writer io.WriteSeeker, v any) error {
	m := marshaller{writer}
	return m.marshal(v)
}
