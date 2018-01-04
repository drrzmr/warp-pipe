package reader

import "io"

// Config struct
type Config struct {
	InputStream  io.Reader
	OutputStream io.WriteCloser
}
