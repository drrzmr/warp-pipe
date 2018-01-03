package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/snippet/reader"
)

// Config main struct
type Config struct {
	Reader   reader.Config
	Log      log.Config
	Postgres postgres.Config
}

const (
	// AppVersion app version
	AppVersion = "0.0.1"
	// AppName app name
	AppName = "warp-pipe"
	// AppShortDescription command line short description
	AppShortDescription = "Golang tools to handle postgres logical replication slots"
	// EnvPrefix prefix for app env vars
	EnvPrefix = "WP"
	// FileType config file type
	FileType = "yaml"
)

var (
	// FilePath for config files
	FilePath = [...]string{
		fmt.Sprintf("/etc/%s", AppName),
		fmt.Sprintf("$HOME/.config/%s", AppName),
		".",
	}
)

// Default config
var Default = &Config{
	Reader: reader.Config{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
	},
	Log: log.Config{
		Stdout: "stdout",
		Stderr: "stderr",
	},
	Postgres: postgres.Config{
		Host:     "postgres",
		Port:     postgres.DefaultPort,
		User:     postgres.DefaultUser,
		Database: "postgres",
		Password: "none",
		Replicate: postgres.ReplicateConfig{
			Slot:   "warp_pipe",
			Plugin: "test_decoding",

			IgnoreDuplicateObjectError: true,
		},
		SQL: postgres.SQLConfig{
			Driver:                   "pgx",
			ConnectTimeout:           10 * time.Second,
			CreateDatabaseIfNotExist: true,
		},
		Streaming: postgres.StreamingReplicateProtocolConfig{
			SendStandByStatusPeriod: 5 * time.Second,
			WaitMessageTimeout:      5 * time.Second,
		},
	},
}

// New the f* 12factor config
func New() (conf *Config) {

	var err error

	// environment variables setup
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType(FileType)

	// template, preload default config to viper
	buf, err := yaml.Marshal(Default)
	no(err)
	err = viper.ReadConfig(bytes.NewBuffer(buf))
	no(err)

	// handle config paths
	viper.SetConfigName(AppName)
	viper.SetConfigType(FileType)
	for _, path := range FilePath {
		viper.AddConfigPath(path)
	}
	// merge configs
	_ = viper.MergeInConfig()

	// dump config to Config struct
	conf = &Config{}
	err = viper.Unmarshal(conf)
	no(err)

	return conf
}

// Uint16 config helper
func Uint16(dst *uint16, src, defaultValue uint16) {
	if src != defaultValue {
		*dst = src
	}
}

// String config helper
func String(dst *string, src, defaultValue string) {
	if src != defaultValue {
		*dst = src
	}
}

func no(err error) {
	if err != nil {
		panic(err)
	}
}
