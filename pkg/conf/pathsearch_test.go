// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_candidateFiles(t *testing.T) {
	gt := NewGomegaWithT(t)

	wd, err := os.Getwd()
	gt.Expect(err).NotTo(HaveOccurred())
	confDir, err := os.UserConfigDir()
	gt.Expect(err).NotTo(HaveOccurred())
	homeDir, err := os.UserHomeDir()
	gt.Expect(err).NotTo(HaveOccurred())

	paths, err := candidateFiles("franklin")
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(paths[0]).To(Equal(filepath.Join(wd, "franklin.yml")))
	gt.Expect(paths[1]).To(Equal(filepath.Join(wd, "franklin.yaml")))
	gt.Expect(paths[2]).To(Equal(filepath.Join(confDir, "franklin", "franklin.yml")))
	gt.Expect(paths[3]).To(Equal(filepath.Join(confDir, "franklin", "franklin.yaml")))
	if filepath.Clean(confDir) != filepath.Join(homeDir, ".config") {
		gt.Expect(paths).To(HaveLen(6))
		gt.Expect(paths[4]).To(Equal(filepath.Join(os.Getenv("HOME"), ".config", "franklin", "franklin.yml")))
		gt.Expect(paths[5]).To(Equal(filepath.Join(os.Getenv("HOME"), ".config", "franklin", "franklin.yaml")))
	}
}

func TestFile(t *testing.T) {
	isHelperProcess := os.Getenv("BATIK_CONFIG_TEST_HELPER_PROCESS") == "1"
	if !isHelperProcess {
		cmd := &exec.Cmd{
			Path: os.Args[0],
			Args: []string{
				os.Args[0],
				"-test.run=" + t.Name(),
				"-test.v=" + strconv.FormatBool(testing.Verbose()),
			},
			Dir:    "testdata",
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Env:    append(os.Environ(), "BATIK_CONFIG_TEST_HELPER_PROCESS=1"),
		}

		err := cmd.Run()
		NewGomegaWithT(t).Expect(err).NotTo(HaveOccurred())
	}

	t.Run("missing", func(t *testing.T) {
		if isHelperProcess {
			t.SkipNow()
		}
		gt := NewGomegaWithT(t)
		path, err := File("missing-config-file-stem")
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(path).To(BeEmpty())
	})

	t.Run("existing", func(t *testing.T) {
		if !isHelperProcess {
			t.SkipNow()
		}

		gt := NewGomegaWithT(t)
		wd, err := os.Getwd()
		gt.Expect(err).NotTo(HaveOccurred())

		path, err := File("batik")
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(path).To(Equal(filepath.Join(wd, "batik.yaml")))
	})
}
