// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import "fmt"

type EnvMap map[string]string

func (m EnvMap) Getenv(key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("$%s is not defined", key)
	}

	return v, nil
}
