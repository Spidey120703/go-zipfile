package serial

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type options struct {
	Length    string `tag:"len"`
	Size      string `tag:"size"`
	Condition string `tag:"condition"`
	Prefix    string `tag:"prefix"`
}

func (o *options) CheckPrefix(r io.ReadSeeker) bool {
	if len(o.Prefix) == 0 {
		return false
	}
	prefix := []byte(o.Prefix)
	if prefix[0] != '\'' || prefix[len(prefix)-1] != '\'' {
		return false
	}
	prefix = prefix[1 : len(prefix)-1]

	data := make([]byte, len(prefix))
	read, err := r.Read(data)
	defer func() { _, _ = r.Seek(-int64(read), io.SeekCurrent) }()

	if err != nil {
		return false
	}
	return bytes.Equal(prefix, data)
}

func getValue(v any, fieldPath string) (reflect.Value, bool) {
	typ := reflect.TypeOf(v).Elem()
	val := reflect.ValueOf(v).Elem()
	fieldNames := strings.Split(fieldPath, ".")

	var found bool
	var field reflect.StructField
	for i, fieldName := range fieldNames {
		if field, found = typ.FieldByName(fieldName); !found {
			return reflect.Value{}, false
		}
		if i == len(fieldNames)-1 {
			return val.Field(field.Index[0]), true
		}
		typ = field.Type
		val = val.FieldByName(fieldName)
	}
	return reflect.Value{}, false
}

func (o *options) GetLength(v any) int {
	if len(o.Length) == 0 {
		return 0
	}

	if value, found := getValue(v, o.Length); found {
		if value.CanUint() {
			return int(value.Uint())
		} else {
			return int(value.Int())
		}
	}
	return 0
}

var funcRegExp = regexp.MustCompile(`([a-z]+)\((.+?), *(.+?)\)`)

func (o *options) GetConditionalResult(v any) bool {
	if len(o.Condition) == 0 {
		return true
	}
	if o.Condition == "false" {
		return false
	}
	if o.Condition == "true" {
		return true
	}
	tokens := funcRegExp.FindStringSubmatch(o.Condition)
	if len(tokens) != 3 {
		return false
	}
	switch tokens[0] {
	case "bit":
		value, found := getValue(v, tokens[1])
		if !found {
			return false
		}
		var flags uint64
		if value.CanUint() {
			flags = value.Uint()
		} else {
			flags = uint64(value.Int())
		}

		imm, err := strconv.Atoi(tokens[2])
		if err != nil {
			return false
		}

		return flags&(1<<imm) != 0
	default:
		return false
	}
}

func parseTag(tag string) (opts options) {
	typ := reflect.TypeOf(opts)
	val := reflect.ValueOf(&opts).Elem()
	tags := strings.Split(tag, ",")
	for _, t := range tags {
		eq := strings.SplitN(t, "=", 2)
		for i := range typ.NumField() {
			if typ.Field(i).Tag.Get("tag") == eq[0] {
				val.Field(i).SetString(eq[1])
			}
		}
	}
	return
}

type unmarshaler struct {
	r io.ReadSeeker
}

func (u *unmarshaler) unmarshal(v, parent any, tag string) (err error) {
	if v == nil {
		return fmt.Errorf("cannot unmarshal nil pointer")
	}

	if deserializable, ok := v.(IDeserializable); ok {
		return deserializable.Unmarshal(u.r)
	}

	varType := reflect.TypeOf(v)
	if varType.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot unmarshal non-pointer value")
	}
	varType = varType.Elem()

	value := reflect.ValueOf(v).Elem()

	switch varType.Kind() {
	case reflect.Struct:
		for i := range varType.NumField() {
			//println(varType.Field(i).Name, varType.Field(i).Tag.Get("serial"))
			if err = u.unmarshal(value.Field(i).Addr().Interface(), v, varType.Field(i).Tag.Get("serial")); err != nil {
				return
			}
		}
	case reflect.Pointer:
		if value.IsNil() {
			opts := parseTag(tag)
			if !opts.GetConditionalResult(parent) {
				return
			}
			value.Set(reflect.New(varType.Elem()))
		}
		return u.unmarshal(value.Interface(), v, "")
	case reflect.Slice:
		opts := parseTag(tag)
		length := opts.GetLength(parent)
		value.Set(reflect.MakeSlice(value.Type(), length, length))
		if value.Len() == 0 {
			for {
				if !opts.CheckPrefix(u.r) {
					break
				}
				elem := reflect.New(value.Type().Elem())
				if err = u.unmarshal(elem.Interface(), v, tag); err != nil {
					return
				}
				value.Set(reflect.Append(value, elem.Elem()))
			}
			return
		}
		fallthrough
	case reflect.Array:
		for i := range value.Len() {
			if err = u.unmarshal(value.Index(i).Addr().Interface(), v, tag); err != nil {
				return
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return binary.Read(u.r, binary.LittleEndian, v)
	default:
		return fmt.Errorf("unsupported type %v", varType)
	}

	return
}

func Unmarshal(reader io.ReadSeeker, v any) error {
	u := unmarshaler{reader}
	return u.unmarshal(v, nil, "")
}
