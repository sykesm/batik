// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Options is a placeholder for configuration structures and command line
// flags.  It's likely that the various options will be moved out of this
// package to live closer to the items that are configured.
//
// More succinctly, this is an experiment to see how component configuration
// may work. If it doesn't, we want to throw it away quickly.
package options

// Configuration Model
// ===================
//
// For each configurable element of Batik, a structure has been defined with
// exported fields. Each of these fields should be documented such that a user
// or developer can understand what the field controls.
//
// In addition to field level documentation, the fields should have be
// annotated with yaml tags that determine the key that can be used in
// configuration files to set the configuration value.
//
// Each configuration structure should also be associated with an exported
// function named the same as the structure with a Defaults suffix. This
// function is used to create a configuration instance with default values
// populated.
//
// The configuration structures should also implement two methods:
//
//   ApplyDefaults() error
//   Flags(commandName string) []cli.Flag
//
// The role of ApplyDefaults is to apply default values to fields that are
// missing a value. This allows us to construct instances of configuration that
// can be used directly by the runtime if populated from a sparse configuraiton
// file or when we add new fields.
//
// The role of Flags is to expose configuration elements as flags on CLI
// commands. The full name of the command and subcommands are provided as
// context. This context enables the same configuration model to be used on
// different commands with possibly different flag names. The Flags method will
// be called after applying defaults.
