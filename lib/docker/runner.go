package docker

import (
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// Runner object
type Runner struct {
	config Config
	client *client.Client
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

// Start create, start and wait for container
func (runner *Runner) Start() (err error) {
	if runner.isStarted() {
		return nil
	}

	if err = runner.newClient(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (runner *Runner) isStarted() (started bool) {

	return runner.client != nil
}

func (runner *Runner) newClient() (err error) {

	runner.client, err = client.NewEnvClient()
	return errors.Wrapf(err, "Could not create docker client, image: %s, name: %s",
		runner.config.ImageName(),
		runner.config.ContainerName,
	)
}
