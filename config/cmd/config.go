package cmd

// DumpConfig config parameters to dump command
type DumpConfig struct {
	Stdout         string
	InputNamedPipe string
}

// Config command config object
type Config struct {
	Dump DumpConfig
}
