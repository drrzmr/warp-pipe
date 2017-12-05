package docker

// Runner object
type Runner struct {
	config Config
}

// NewRunner return new Runner object
func NewRunner(config Config) *Runner {

	return &Runner{
		config: config,
	}
}

// Config return runner config object
func (runner *Runner) Config() *Config {
	return &runner.config
}
