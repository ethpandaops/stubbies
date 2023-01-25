package server

import (
	"github.com/ethpandaops/stubbies/pkg/execution"
)

type Config struct {
	LoggingLevel string `yaml:"logging" default:"info"`
	Addr         string `yaml:"addr" default:":8551"`
	MetricsAddr  string `yaml:"metricsAddr" default:":9090"`

	Execution execution.Config `yaml:"execution"`
}

func (c *Config) Validate() error {
	return nil
}
