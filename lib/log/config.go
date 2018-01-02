package log

// Config log object
type Config struct {
	Stdout string
	Stderr string
}

// Default log config
var Default = Config{
	Stdout: "stdout",
	Stderr: "stderr",
}

// Test log config
var Test = Default
