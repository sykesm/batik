// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestMapLookuper(t *testing.T) {
	gt := NewGomegaWithT(t)

	ml := MapLookuper(nil)
	v, ok := ml.Lookup("key1")
	gt.Expect(ok).To(BeFalse())
	gt.Expect(v).To(BeEmpty())

	ml = MapLookuper(map[string]string{"key": "value"})
	v, ok = ml.Lookup("key")
	gt.Expect(ok).To(BeTrue())
	gt.Expect(v).To(Equal("value"))

	v, ok = ml.Lookup("missing")
	gt.Expect(ok).To(BeFalse())
	gt.Expect(v).To(BeEmpty())
}

func TestEnvironLookuper(t *testing.T) {
	testEnviron := []string{
		"LOOKUPTEST_KEY_ONE=environment_value_one",
		"lookuptest_key_two=value_two",
	}

	if os.Getenv("LOOKUPTEST_WANT_HELPER_PROCESS") != "1" {
		// Explicitly filter any environment starting with LOOKUPTEST or lookuptest
		env := os.Environ()[:0]
		for _, e := range env {
			if !strings.HasPrefix(e, "LOOKUPTEST_") && !strings.HasPrefix(e, "lookuptest_") {
				env = append(env, e)
			}
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestEnvironLookuper")
		cmd.Env = append(env, "LOOKUPTEST_WANT_HELPER_PROCESS=1")
		cmd.Env = append(cmd.Env, testEnviron...)
		out, err := cmd.CombinedOutput()
		NewGomegaWithT(t).Expect(err).NotTo(HaveOccurred(), string(out))
		return
	}

	for _, kv := range testEnviron {
		t.Run(kv, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tokens := strings.SplitN(kv, "=", 2)
			gt.Expect(tokens).To(HaveLen(2))
			key, val := tokens[0], tokens[1]

			v, ok := EnvironLookuper().Lookup(key)
			gt.Expect(ok).To(BeTrue())
			gt.Expect(v).To(Equal(val))
		})
	}

	t.Run("lookupenv_missing", func(t *testing.T) {
		_, ok := EnvironLookuper().Lookup("lookupenv_missing")
		NewGomegaWithT(t).Expect(ok).To(BeFalse())
	})
}
