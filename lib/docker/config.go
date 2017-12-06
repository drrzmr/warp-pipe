package docker

import (
	"fmt"
	"time"
)

const (
	minWaitTimeout = 10 * time.Second
	// DefaultWaitTimeout default timeout
	DefaultWaitTimeout = 5 * minWaitTimeout
)

// Config object
type Config struct {
	Image         string
	URL           string
	Tag           string
	ContainerName string
	WaitTimeout   time.Duration
}

// ImageName return image name from config
func (c *Config) ImageName() string {

	var tag string

	if len(c.Tag) == 0 {
		tag = "latest"
	} else {
		tag = c.Tag
	}

	return fmt.Sprintf("%s/%s:%s", c.URL, c.Image, tag)
}
