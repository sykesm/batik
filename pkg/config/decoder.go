// Copyright (c) 2015-2019 Carlos Alexandro Becker
// The MIT License (MIT)

package config

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
	// Struct to Parse
	ErrNotAStructPtr = errors.New("decode: expected a pointer to a Struct")

	defaultBuiltInParsers = map[reflect.Kind]ParserFunc{
		reflect.Bool: func(v string) (interface{}, error) {
			return strconv.ParseBool(v)
		},
		reflect.String: func(v string) (interface{}, error) {
			return v, nil
		},
		reflect.Int: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 32)
			return int(i), err
		},
		reflect.Int16: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 16)
			return int16(i), err
		},
		reflect.Int32: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 32)
			return int32(i), err
		},
		reflect.Int64: func(v string) (interface{}, error) {
			return strconv.ParseInt(v, 10, 64)
		},
		reflect.Int8: func(v string) (interface{}, error) {
			i, err := strconv.ParseInt(v, 10, 8)
			return int8(i), err
		},
		reflect.Uint: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 32)
			return uint(i), err
		},
		reflect.Uint16: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 16)
			return uint16(i), err
		},
		reflect.Uint32: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 32)
			return uint32(i), err
		},
		reflect.Uint64: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 64)
			return i, err
		},
		reflect.Uint8: func(v string) (interface{}, error) {
			i, err := strconv.ParseUint(v, 10, 8)
			return uint8(i), err
		},
		reflect.Float64: func(v string) (interface{}, error) {
			return strconv.ParseFloat(v, 64)
		},
		reflect.Float32: func(v string) (interface{}, error) {
			f, err := strconv.ParseFloat(v, 32)
			return float32(f), err
		},
	}

	defaultTypeParsers = map[reflect.Type]ParserFunc{
		reflect.TypeOf(url.URL{}): func(v string) (interface{}, error) {
			u, err := url.Parse(v)
			if err != nil {
				return nil, errors.Wrap(err, "unable to parse URL")
			}
			return *u, nil
		},
		reflect.TypeOf(time.Nanosecond): func(v string) (interface{}, error) {
			s, err := time.ParseDuration(v)
			if err != nil {
				return nil, errors.Wrap(err, "unable to parse duration")
			}
			return s, err
		},
	}
)

type Decoder struct {
	lookuper   Lookuper
	parseTag   string
	defaultTag string
}

// ParserFunc defines the signature of a function that can be used within `CustomParsers`
type ParserFunc func(v string) (interface{}, error)

func (d Decoder) Parse(v interface{}) error {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}

	return d.doParse(ref)
}

func (d Decoder) doParse(ref reflect.Value) error {
	var refType = ref.Type()

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		if !refField.CanSet() {
			continue
		}
		if reflect.Ptr == refField.Kind() && !refField.IsNil() {
			err := d.Parse(refField.Interface())
			if err != nil {
				return err
			}
			continue
		}
		if reflect.Struct == refField.Kind() && refField.CanAddr() && refField.Type().Name() == "" {
			err := d.Parse(refField.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}
		refTypeField := refType.Field(i)
		parseTagFound := true
		value, ok, err := d.get(refTypeField)
		if err != nil {
			return err
		}
		if !ok {
			value = refTypeField.Tag.Get(d.defaultTag)
			parseTagFound = false
		}
		if value == "" {
			if reflect.Struct == refField.Kind() {
				if err := d.doParse(refField); err != nil {
					return err
				}
			}
			continue
		}
		// If the field is already set and the env var is not set don't override with default
		if !refField.IsZero() && !parseTagFound {
			continue
		}
		if err := d.set(refField, refTypeField, value); err != nil {
			return err
		}
	}
	return nil
}

func (d Decoder) get(field reflect.StructField) (val string, found bool, err error) {
	key, opts := parseKeyAndOptions(field.Tag.Get(d.parseTag))

	for _, opt := range opts {
		switch opt {
		case "":
			break
		default:
			return "", false, errors.Errorf("decode: tag option %q not supported", opt)
		}
	}

	val, ok := d.lookuper.Lookup(key)
	return val, ok, nil
}

func (d Decoder) set(field reflect.Value, sf reflect.StructField, value string) error {
	if field.Kind() == reflect.Slice {
		return handleSlice(field, value, sf)
	}

	var tm = asTextUnmarshaler(field)
	if tm != nil {
		var err = tm.UnmarshalText([]byte(value))
		return newParseError(sf, err)
	}

	var typee = sf.Type
	var fieldee = field
	if typee.Kind() == reflect.Ptr {
		typee = typee.Elem()
		fieldee = field.Elem()
	}

	parserFunc, ok := defaultTypeParsers[typee]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return newParseError(sf, err)
		}

		fieldee.Set(reflect.ValueOf(val))
		return nil
	}

	parserFunc, ok = defaultBuiltInParsers[typee.Kind()]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return newParseError(sf, err)
		}

		fieldee.Set(reflect.ValueOf(val).Convert(typee))
		return nil
	}

	return newNoParserError(sf)
}

// parseKeyAndOptions splits the tag's key into the expected key and desired options, if any.
func parseKeyAndOptions(key string) (string, []string) {
	opts := strings.Split(key, ",")
	return opts[0], opts[1:]
}

func handleSlice(field reflect.Value, value string, sf reflect.StructField) error {
	separator := ","
	var parts = strings.Split(value, separator)

	var typee = sf.Type.Elem()
	if typee.Kind() == reflect.Ptr {
		typee = typee.Elem()
	}

	if _, ok := reflect.New(typee).Interface().(encoding.TextUnmarshaler); ok {
		return parseTextUnmarshalers(field, parts, sf)
	}

	parserFunc, ok := defaultTypeParsers[typee]
	if !ok {
		parserFunc, ok = defaultBuiltInParsers[typee.Kind()]
		if !ok {
			return newNoParserError(sf)
		}
	}

	var result = reflect.MakeSlice(sf.Type, 0, len(parts))
	for _, part := range parts {
		r, err := parserFunc(part)
		if err != nil {
			return newParseError(sf, err)
		}
		var v = reflect.ValueOf(r).Convert(typee)
		if sf.Type.Elem().Kind() == reflect.Ptr {
			v = reflect.New(typee)
			v.Elem().Set(reflect.ValueOf(r).Convert(typee))
		}
		result = reflect.Append(result, v)
	}
	field.Set(result)
	return nil
}

func asTextUnmarshaler(field reflect.Value) encoding.TextUnmarshaler {
	if reflect.Ptr == field.Kind() {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	} else if field.CanAddr() {
		field = field.Addr()
	}

	tm, ok := field.Interface().(encoding.TextUnmarshaler)
	if !ok {
		return nil
	}
	return tm
}

func parseTextUnmarshalers(field reflect.Value, data []string, sf reflect.StructField) error {
	s := len(data)
	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(elemType), s, s)
	for i, v := range data {
		sv := slice.Index(i)
		kind := sv.Kind()
		if kind == reflect.Ptr {
			sv = reflect.New(elemType.Elem())
		} else {
			sv = sv.Addr()
		}
		tm := sv.Interface().(encoding.TextUnmarshaler)
		if err := tm.UnmarshalText([]byte(v)); err != nil {
			return newParseError(sf, err)
		}
		if kind == reflect.Ptr {
			slice.Index(i).Set(sv)
		}
	}

	field.Set(slice)

	return nil
}

func newParseError(sf reflect.StructField, err error) error {
	if err == nil {
		return nil
	}
	return parseError{
		sf:  sf,
		err: err,
	}
}

type parseError struct {
	sf  reflect.StructField
	err error
}

func (e parseError) Error() string {
	return fmt.Sprintf(`decode: parse error on field "%s" of type "%s": %v`, e.sf.Name, e.sf.Type, e.err)
}

func newNoParserError(sf reflect.StructField) error {
	return errors.Errorf(`decode: no parser found for field "%s" of type "%s"`, sf.Name, sf.Type)
}
