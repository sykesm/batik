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

- [ ] Remove format strings and fabenc
- [ ] Write a colorized post processor for logfmt (see humanlog patterns)
- [ ] Process logging config options as flags (eg. color)
- [ ] Introduce terminal detection for color processing
- [ ] Helpers (eg IsColorEnabled) for config object that wraps logger
- [ ] Reduce error processing in logger creation path (handle errors higher up in app)

# DB

- [ ] Configure and create a persistent DB
- [ ] Introduce commands for interacting with db directly (get/put)
- [ ] Introduce command to dump keys

# Errors

- [ ] Evaluate consistent treatment of errors and Wrap vs. WithMessage

## Discussion Topics

- [ ] The best way to organize this list and regularly groom it
- [x] How to track completed tasks (leave in place or move?)
- [x] How to handle logging and error streams (one vs two loggers?)
