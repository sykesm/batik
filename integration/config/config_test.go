// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v3"

	"github.com/sykesm/batik/pkg/options"
)

var _ = Describe("Command Line Configuration", func() {
	var tempdir string

	BeforeEach(func() {
		var err error
		tempdir, err = ioutil.TempDir("", "config")
		Expect(err).NotTo(HaveOccurred())
		tempdir = resolveSymlinks(tempdir)
	})

	AfterEach(func() {
		if tempdir != "" {
			os.RemoveAll(tempdir)
		}
	})

	Context("without a configuration file", func() {
		It("uses default values", func() {
			cmd := exec.Command(batikPath, "--show-config")
			cmd.Dir = tempdir

			config := decodeShowConfigOutput(cmd)
			Expect(config).To(Equal(options.BatikDefaults()))
		})
	})

	Context("with an empty configuration file", func() {
		BeforeEach(func() {
			err := ioutil.WriteFile(filepath.Join(tempdir, "batik.yaml"), []byte("---\n"), 0o644)
			Expect(err).NotTo(HaveOccurred())
		})

		It("uses default values", func() {
			cmd := exec.Command(batikPath, "--show-config")
			cmd.Dir = tempdir

			expected := options.BatikDefaults()
			expected.DataDir = filepath.Join(tempdir, "data")
			expected.Server.TLS.CertsDir = filepath.Join(tempdir, "tls-certs")

			config := decodeShowConfigOutput(cmd)
			Expect(config).To(Equal(expected))
		})

		It("honors global flags", func() {
			cmd := exec.Command(
				batikPath, "--show-config",
				"--color", "yes",
				"--log-format", "json",
				"--log-spec", "error",
			)
			cmd.Dir = tempdir

			expected := options.BatikDefaults()
			expected.DataDir = filepath.Join(tempdir, "data")
			expected.Server.TLS.CertsDir = filepath.Join(tempdir, "tls-certs")
			expected.Logging.Color = "yes"
			expected.Logging.Format = "json"
			expected.Logging.LogSpec = "error"

			config := decodeShowConfigOutput(cmd)
			Expect(config).To(Equal(expected))
		})

		It("honors data-dir as a relative path", func() {
			cmd := exec.Command(
				batikPath, "--show-config",
				"--data-dir", filepath.Join("custom-data-dir"),
			)
			cmd.Dir = tempdir

			expected := options.BatikDefaults()
			expected.DataDir = filepath.Join(tempdir, "custom-data-dir")
			expected.Server.TLS.CertsDir = filepath.Join(tempdir, "tls-certs")

			config := decodeShowConfigOutput(cmd)
			Expect(config).To(Equal(expected))
		})

		It("honors tls-certs-dir as a relative path", func() {
			cmd := exec.Command(
				batikPath, "--show-config",
				"start",
				"--tls-certs-dir", filepath.Join(tempdir, "some/custom-tls-dir"),
			)
			cmd.Dir = tempdir

			expected := options.BatikDefaults()
			expected.DataDir = filepath.Join(tempdir, "data")
			expected.Server.TLS.CertsDir = filepath.Join(tempdir, "some/custom-tls-dir")

			config := decodeShowConfigOutput(cmd)
			Expect(config).To(Equal(expected))
		})

		It("honors server flags", func() {
			cmd := exec.Command(
				batikPath, "--show-config",
				"start",
				"--grpc-conn-timeout", "5s",
				"--grpc-listen-address", "grpc:123",
				"--grpc-max-recv-message-size", "999",
				"--grpc-max-send-message-size", "888",
				"--http-listen-address", "http:http",
				"--http-read-timeout", "90s",
				"--http-read-header-timeout", "1s",
				"--http-write-timeout", "1m",
				"--http-idle-timeout", "5m",
				"--tls-cert-file", "my-tls.pem",
				"--tls-private-key-file", "my-private-tls.key",
			)
			cmd.Dir = tempdir

			expected := options.BatikDefaults()
			expected.DataDir = filepath.Join(tempdir, "data")
			expected.Server.TLS = options.ServerTLS{
				ServerCert: options.CertKeyPair{
					CertFile: "my-tls.pem",
					KeyFile:  "my-private-tls.key",
				},
				CertsDir: filepath.Join(tempdir, "tls-certs"),
			}
			expected.Server.GRPC = options.GRPCServer{
				GRPC: options.GRPC{
					MaxRecvMessageSize: 999,
					MaxSendMessageSize: 888,
				},
				ConnTimeout:   5 * time.Second,
				ListenAddress: "grpc:123",
			}
			expected.Server.HTTP = options.HTTPServer{
				ListenAddress:     "http:http",
				ReadTimeout:       90 * time.Second,
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      60 * time.Second,
				IdleTimeout:       5 * time.Minute,
			}

			config := decodeShowConfigOutput(cmd)
			Expect(config).To(Equal(expected))
		})
	})

	Context("with a sparse configuration file", func() {
		It("honors data_dir in the configuration file", func() {
			config := &options.Batik{
				DataDir: "/absolute/path",
			}
			writeConfigurationFile(filepath.Join(tempdir, "batik.yaml"), config)

			cmd := exec.Command(batikPath, "--show-config")
			cmd.Dir = tempdir

			config.ApplyDefaults()
			config.Server.TLS.CertsDir = filepath.Join(tempdir, "tls-certs")

			active := decodeShowConfigOutput(cmd)
			Expect(active).To(Equal(config))
		})

		// This tests demonstrate that global flags do not override values
		// from a configuration file.
		PIt("global flags are prioritized over the configuration file", func() {
			config := &options.Batik{
				DataDir: "/absolute/path",
			}
			writeConfigurationFile(filepath.Join(tempdir, "batik.yaml"), config)

			cmd := exec.Command(
				batikPath,
				"--show-config",
				"--data-dir", "another/data/dir",
			)
			cmd.Dir = tempdir

			config.ApplyDefaults()
			config.DataDir = filepath.Join(tempdir, "another/data/dir")
			config.Server.TLS.CertsDir = filepath.Join(tempdir, "tls-certs")

			active := decodeShowConfigOutput(cmd)
			Expect(active).To(Equal(config))
		})

		It("command flags are prioritized over the configuration file", func() {
			config := &options.Batik{
				Server: options.Server{
					HTTP: options.HTTPServer{
						ListenAddress: "http:http",
					},
				},
			}
			writeConfigurationFile(filepath.Join(tempdir, "batik.yaml"), config)

			cmd := exec.Command(
				batikPath,
				"--show-config",
				"start",
				"--http-listen-address", "https:https",
			)
			cmd.Dir = tempdir

			config.ApplyDefaults()
			config.DataDir = filepath.Join(tempdir, "data")
			config.Server.HTTP.ListenAddress = "https:https"
			config.Server.TLS.CertsDir = filepath.Join(tempdir, "tls-certs")

			active := decodeShowConfigOutput(cmd)
			Expect(active).To(Equal(config))
		})
	})
})

func decodeShowConfigOutput(cmd *exec.Cmd) *options.Batik {
	var config options.Batik
	sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(sess, time.Minute).Should(gexec.Exit(0))

	err = yaml.NewDecoder(bytes.NewBuffer(sess.Out.Contents())).Decode(&config)
	Expect(err).NotTo(HaveOccurred())
	return &config
}

func writeConfigurationFile(path string, config *options.Batik) {
	f, err := os.Create(path)
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(config)
	Expect(err).NotTo(HaveOccurred())
}

// macOS makes /tmp a symlink to /private/tmp
func resolveSymlinks(p string) string {
	resolved, err := filepath.EvalSymlinks(p)
	Expect(err).NotTo(HaveOccurred())
	return resolved
}
