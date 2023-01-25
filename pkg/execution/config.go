package execution

type Config struct {
	ChainID                 string `yaml:"chainID" default:"0x1"`
	TerminalTotalDifficulty string `yaml:"terminalTotalDifficulty" default:"0x0"`
	TerminalBlockHash       string `yaml:"terminalBlockHash" default:"0x0000000000000000000000000000000000000000000000000000000000000000"`
	TerminalBlockNumber     string `yaml:"terminalBlockNumber" default:"0x0"`
}

func (c *Config) Validate() error {
	return nil
}
