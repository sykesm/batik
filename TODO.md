# Living TODO List

This document contains a set of stories and work items for Batik.

## Configuration

- [ ] Investigate the `altsrc` package for loading configuration data from different sources
- [ ] Investigate "defaulters", "option" structs, and the "AddFlags" patterns for setting up and overriding configuration options for each command
- [ ] Investigate doc generation from the cli package

## Services

- [ ] Adopt ifrit (or its patterns) for managing multiple independent processes in the server
- [x] Introduce zap logger for logging

## Discussion Topics

- [ ] The best way to organize this list and regularly groom it
- [x] How to track completed tasks (leave in place or move?)
- [x] How to handle logging and error streams (one vs two loggers?)
