# Living TODO List

This document contains a set of stories and work items for Batik.

## Configuration

- [x] Investigate the `altsrc` package for loading configuration data from different sources
- [x] Investigate "defaulters", "option" structs, and the "AddFlags" patterns
      for setting up and overriding configuration options for each command
- [ ] Create runtime TLS artifacts from configuration
- [ ] Investigate doc generation from the cli package
  - [x] CLI documentation in markdown
        > cmd/docgen generates something we can start from
  - [x] CLI man page
        > cmd/docgen generates something we can start from
  - [ ] Example YAML configuration from defaults + tags
        > I started looking into this and it's simply not worth the effort at
        > this time. The general idea is to take the default configuration
        > object, use reflection to get the types referenced by the object,
        > then use x/tools/packages and go/ast to get the syntax tree for the
        > structure definitions. From this tree we can extract comments. This
        > is actually pretty straightforward and a small amount of code.
        >
        > For the yaml, we can marshal the default configuration to a yaml.v3
        > Node. The node represents the syntax tree of the yaml document. For
        > simple documents, the structure is reasonable.
        >
        > With the default config object, the marshaled yaml node, and the AST
        > info, the idea is to walk the yaml document and add comments to
        > nodes in the tree from the comments in the AST. If an example tag is
        > present on the field, replace node values with example values. This
        > seems pretty simple but it can get tricky if we have to deal with
        > omitted fields, inline fields, anchors, or aliases.
        >
        > There are also bugs in the yaml package related to where comments
        > are placed in the tree so we would probably end up fighting those.

## Services

- [ ] Adopt ifrit (or its patterns) for managing multiple independent processes in the server
  - [x] gRPC
- [x] Introduce zap logger for logging

## Logging

- [x] Remove format strings and fabenc
- [x] Write a colorized post processor for logfmt (see humanlog patterns)
- [x] Process logging config options as flags (eg. color)
- [x] Introduce terminal detection for color processing
- [ ] Helpers (eg IsColorEnabled) for config object that wraps logger
- [x] Reduce error processing in logger creation path (handle errors higher up in app)

# DB

- [x] Configure and create a persistent DB
- [x] Introduce commands for interacting with db directly (get/put)
      > `db get` and `db put` subcommands available in both interactive and non interactive.
      > `db get` expects a hex encoded string key and returns a go hex.Dump of the stored
      > value.
      > `db put` expects hex encoded strings for both the key and value and stores the
      > decoded value at the decoded key.
- [x] Introduce command to dump keys
      > `db keys` returns a list of hex encoded string keys over the entire db.
- [ ] Create and manage data access errors consistently. In particular, when a
      value does not exist, we need an error that clearly indicates that. When
      a value cannot be unmarshaled into a message, we need a different error.
      The first case is likely normal while the second should be fatal for us.

# Errors

- [ ] Evaluate consistent treatment of errors and Wrap vs. WithMessage

## Discussion Topics

- [ ] The best way to organize this list and regularly groom it
- [x] How to track completed tasks (leave in place or move?)
- [x] How to handle logging and error streams (one vs two loggers?)
