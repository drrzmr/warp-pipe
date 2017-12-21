package file

type config struct {
	InputExtension  string
	OutputExtension string
	JSONExtension   string
}

// Config object
var Config = config{
	InputExtension:  "in",
	OutputExtension: "out",
	JSONExtension:   "json",
}
