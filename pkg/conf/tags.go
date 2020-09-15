// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"
)

type TagResolver struct {
	SourcePath string
}

func (tr *TagResolver) Resolve(input interface{}) error {
	if input == nil {
		return nil
	}

	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() || v.Elem().Kind() != reflect.Struct {
			return nil
		}
	case reflect.Struct:
	default:
		return nil
	}

	return tr.resolve(v)
}

func (tr *TagResolver) resolve(v reflect.Value) error {
	v = reflect.Indirect(v)
	t := v.Type()

	var err error
	for i := 0; i < t.NumField() && err == nil; i++ {
		switch t.Field(i).Type.Kind() {
		case reflect.Ptr, reflect.Interface:
			if !v.Field(i).IsNil() {
				err = tr.resolve(v.Field(i).Elem())
			}
		case reflect.Struct:
			err = tr.resolve(v.Field(i))
		default:
			if tag, ok := t.Field(i).Tag.Lookup("batik"); ok {
				err = tr.resolveTag(tag, v.Field(i))
			}
		}
	}
	return err
}

func (tr *TagResolver) resolveTag(tag string, v reflect.Value) error {
	if !v.CanAddr() {
		return errors.New("field must be addressable")
	}
	switch tag {
	case "relpath":
		return tr.resolveRelpath(v)
	default:
		return errors.Errorf("unknown directive: %s", tag)
	}
}

func (tr *TagResolver) resolveRelpath(v reflect.Value) error {
	if v.Type().Kind() != reflect.String {
		return errors.New("field must be a string")
	}
	path := v.Interface().(string)
	switch {
	case path == "":
	case filepath.IsAbs(path):
	default:
		path = filepath.Join(tr.SourcePath, path)
		v.SetString(path)
	}
	return nil
}
