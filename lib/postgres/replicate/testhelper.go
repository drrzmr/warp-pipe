package replicate

import (
	"testing"
	"time"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/pagarme/warp-pipe/lib/postgres"
)

// CreateTestPostgresConfig test helper
func CreateTestPostgresConfig(t *testing.T) (config postgres.Config) {
	t.Helper()

	return postgres.Config{
		Host:     "none.host",
		Port:     postgres.DefaultPort,
		User:     postgres.DefaultUser,
		Database: "test-replicate",
		Password: "postgres",
		Replicate: postgres.ReplicateConfig{
			Slot:   "test_replicate_slot",
			Plugin: "test_decoding",
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
	}
}

// CreateTestDockerConfig test helper
func CreateTestDockerConfig(t *testing.T) (config docker.Config) {
	t.Helper()

	return docker.Config{
		WaitTimeout: docker.DefaultWaitTimeout,
		URL:         "warp-pipe",
		Image:       "postgres-server",
		Tag:         "9.5.6",
	}
}
