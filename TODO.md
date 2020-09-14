# Living TODO List

This document contains a set of stories and work items for Batik.

## Configuration

- [ ] Investigate the `altsrc` package for loading configuration data from different sources
- [ ] Investigate "defaulters", "option" structs, and the "AddFlags" patterns for setting up and overriding configuration options for each command
  In progress
- [ ] Investigate doc generation from the cli package

## Services

- [ ] Adopt ifrit (or its patterns) for managing multiple independent processes in the server
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
