package namedpipe

import (
	"bufio"
	"context"
	"os"
	"syscall"
)

// NamedPipe object
type NamedPipe struct {
	name string
}

// New namedpipe
func New(filename string) *NamedPipe {

	return &NamedPipe{
		name: filename,
	}
}

// Create pipe file
func (pipe *NamedPipe) Create() (err error) {

	err = syscall.Mkfifo(pipe.name, 0600)
	if err == syscall.EEXIST {
		return nil
	}

	return err
}

// Read return lines from pipe
func (pipe *NamedPipe) Read() (lineList []string, err error) {

	if err = pipe.Create(); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(pipe.name, os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineList = append(lineList, scanner.Text())
	}

	return lineList, scanner.Err()
}

// Loop look for lines on pipe
func (pipe *NamedPipe) Loop(ctx context.Context) (lineCh <-chan string, errCh <-chan error) {

	var (
		lineChannel = make(chan string)
		errChannel  = make(chan error)
	)

	go func(ctx context.Context, lineCh chan<- string, errCh chan<- error) {
		defer close(lineCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				lineList, err := pipe.Read()
				if err != nil {
					errCh <- err
				}

				for _, line := range lineList {
					lineCh <- line
				}
			}

		}
	}(ctx, lineChannel, errChannel)

	return lineChannel, errChannel
}
