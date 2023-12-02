package config

const (
	DebugMode   = "debug"
	TestMode    = "test"
	ReleaseMode = "release"
)

type Config struct {
	Environment string

	ServiceHost string
	HTTPPort    string
	HTTPScheme  string
	Domain      string

	SwaggerTitle   string
	SwaggerVersion string

	InputFilePath string
	RecordingsDir string
}

func Load() Config {
	config := Config{}

	config.Environment = DebugMode

	config.ServiceHost = "localhost"
	config.HTTPPort = "8000"
	config.HTTPScheme = "http"
	config.Domain = "localhost:8000"

	// swagger
	config.SwaggerTitle = "SAS Project"
	config.SwaggerVersion = "release"

	config.InputFilePath = "./input/input_recording.wav"
	config.RecordingsDir = "./recordings"

	return config
}
