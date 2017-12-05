package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/retry"
	"github.com/pagarme/warp-pipe/lib/waitfor"
)

// Runner object
type Runner struct {
	config      Config
	context     context.Context
	containerID string
	json        *types.ContainerJSON
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

	if err = runner.start(); err != nil {
		return errors.WithStack(err)
	}

	if err = runner.waitForRunning(); err != nil {
		return errors.WithStack(err)
	}

	if err = runner.waitForTCPPorts(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Stop stop and remove a started container
func (runner *Runner) Stop() (err error) {

	if runner.isStopped() {
		return nil
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

func (runner *Runner) start() (err error) {

	startOptions := types.ContainerStartOptions{}

	err = runner.client.ContainerStart(runner.context, runner.containerID, startOptions)
	return errors.Wrapf(err, "Could not start docker container, image: %s, name: %s",
		runner.config.ImageName(),
		runner.config.ContainerName,
	)
}

func (runner *Runner) inspect() (json *types.ContainerJSON, err error) {

	var ret types.ContainerJSON

	if ret, err = runner.client.ContainerInspect(runner.context, runner.containerID); err != nil {
		return nil, errors.Wrapf(err, "Could not inspect docker container, image: %s, name: %s",
			runner.config.ImageName(),
			runner.config.ContainerName,
		)
	}

	return &ret, nil
}

func (runner *Runner) waitForRunning() (err error) {

	var (
		json     *types.ContainerJSON
		innerErr error
	)

	err, innerErr = retry.Do(runner.config.WaitTimeout, func() (err error) {

		if json, err = runner.inspect(); err != nil {
			return errors.WithStack(err)
		}

		if !json.State.Running {
			return retry.ErrContinue
		}

		runner.json = json // Store first inspect after running state
		return nil         // Exit with success

	})

	if err != nil {
		return errors.Wrapf(innerErr, "retry end by %s", err.Error())
	}

	return nil
}

func (runner *Runner) waitForTCPPorts() (err error) {

	ipAddress := runner.json.NetworkSettings.IPAddress
	for containerPort := range runner.json.NetworkSettings.Ports {
		if containerPort.Proto() != "tcp" {
			continue
		}
		port := uint16(containerPort.Int())

		if err = waitfor.TCPPort(runner.config.WaitTimeout, ipAddress, port); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (runner *Runner) isStopped() (stopped bool) {

	return runner.client == nil
}
