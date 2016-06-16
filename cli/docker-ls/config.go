package main

import (
	"flag"
)

const (
	OPTION_RECURSION_LEVEL = 1 << iota
	OPTION_STATISTICS
	OPTION_PROGRESS
	OPTION_JSON_OUTPUT
	OPTION_MANIFEST_VERSION
)

const (
	OPTIONS_FULL = 0xFFFF
	OPTIONS_NONE = 0
)

type Config struct {
	recursionLevel  uint
	manifestVersion uint
	statistics      bool
	progress        bool
	jsonOutput      bool
}

func (c *Config) bindToFlags(flags *flag.FlagSet, options uint) {
	if options&OPTION_RECURSION_LEVEL != 0 {
		flags.UintVar(&c.recursionLevel, "level", c.recursionLevel, "level of recursion")
	}

	if options&OPTION_STATISTICS != 0 {
		flags.BoolVar(&c.statistics, "statistics", c.statistics, "show connection statistics")
	}

	if options&OPTION_PROGRESS != 0 {
		flags.BoolVar(&c.progress, "progress-indicator", c.progress, "show progress indicator")
	}

	if options&OPTION_JSON_OUTPUT != 0 {
		flags.BoolVar(&c.jsonOutput, "json", c.jsonOutput, "output JSON instead of YAML")
	}

	if options&OPTION_MANIFEST_VERSION != 0 {
		flags.UintVar(&c.manifestVersion, "manifest-version", c.manifestVersion, "manifest version to request")
	}
}

func newConfig() *Config {
	return &Config{
		recursionLevel:  0,
		statistics:      false,
		progress:        true,
		jsonOutput:      false,
		manifestVersion: 2,
	}
}
