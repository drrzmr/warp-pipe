package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// Runner object
type Runner struct {
	config      Config
	context     context.Context
	containerID string
	client      *client.Client
}

// NewRunner return new Runner object
func NewRunner(config Config) *Runner {

	return &Runner{
		config:  config,
		context: context.Background(),
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

	if err = runner.create(); err != nil {
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

func (runner *Runner) create() (err error) {

	containerConfig := &container.Config{
		Image: runner.config.ImageName(),
	}

	var body container.ContainerCreateCreatedBody

	body, err = runner.client.ContainerCreate(
		runner.context,
		containerConfig,
		nil,
		nil,
		runner.config.ContainerName,
	)

	if err != nil {
		return errors.Wrapf(err, "Could not create docker container, image: %s, name: %s",
			runner.config.ImageName(),
			runner.config.ContainerName,
		)
	}

	runner.containerID = body.ID
	return nil
}
