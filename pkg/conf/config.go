// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

// LoadFile loads a YAML configuration document from the provided path into the
// object referenced by out. After loading the configuration, ApplyDefaults, if
// implemented, will be invoked on out.
//
// Once the configuration has been loaded and the defaults applied, the fields
// of out annotated with the "batik" tag will be recursively processed by an
// instance of the tag resolver setup to use the config source path as the
// relative path root.
func LoadFile(path string, out interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "conf")
	}
	defer f.Close()

	tr := &TagResolver{SourcePath: filepath.Dir(path)}
	return load(f, tr, out)
}

// Load loads a YAML configuration document from teh provided reader into the
// object referenced by out. After loading the configuration, ApplyDefaults, if
// implemented, will be invoked on out.
//
// Once the configuration has been loaded and the defaults applied, the fields
// of out annotated with the "batik" tag will be recursively processed by an
// instance of the tag resolver setup to use the current working directory as
// the relative path root.
func Load(r io.Reader, out interface{}) error {
	return load(r, &TagResolver{}, out)
}

func load(r io.Reader, tr *TagResolver, out interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(out); err != nil {
		return errors.WithStack(err)
	}
	ad, ok := out.(interface{ ApplyDefaults() error })
	if ok {
		if err := ad.ApplyDefaults(); err != nil {
			return err
		}
	}
	return tr.Resolve(out)
}
