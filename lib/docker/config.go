package docker

import "fmt"

// Config object
type Config struct {
	Image         string
	URL           string
	Tag           string
	ContainerName string
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
