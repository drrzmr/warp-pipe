package config

import (
	"io"
	"os"
)

// Defaults
const (
	// AppVersion app version
	AppVersion = "0.0.1"
	// AppName app name
	AppName = "warp-pipe"
	// AppShortDescription command line short description
	AppShortDescription = "Golang tools to handle postgres logical replication slots"
	// ConfigEnvPrefix prefix for app env vars
	ConfigEnvPrefix = "WP"
	// ConfigFileType config filetype
	ConfigFileType = "yaml"
)

// Reader config struct
type Reader struct {
	InputStream  io.Reader
	OutputStream io.WriteCloser
}

// Config main struct
type Config struct {
	Reader
}

// New create a new Config struct
func New() (conf *Config) {
	return &Config{
		Reader{
			InputStream:  os.Stdin,
			OutputStream: os.Stdout,
		},
	}
}
