package collector_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx"

	postgresCollector "github.com/pagarme/warp-pipe/adapter/collector/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	postgresTester "github.com/pagarme/warp-pipe/lib/tester/postgres"
	"github.com/pagarme/warp-pipe/pipeline/collector"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPostgresCollector(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	dockerConfig := replicate.CreateTestDockerConfig(t)
	postgresConfig := replicate.CreateTestPostgresConfig(t)

	normal, deferFn := postgresTester.DockerRun(t, dockerConfig, &postgresConfig)
	defer deferFn()

	normalDB := normal.DB()
	require.NotNil(t, normalDB)

	_, err := normalDB.Exec(`
CREATE TABLE test
(
	id   SERIAL,
	name VARCHAR(30),
	ts   TIMESTAMP NOT NULL,
	PRIMARY KEY (id)
);`)
	require.NoError(t, err)

	t.Run("BuildStage", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		publishCh, offsetCh, err := collector.Run(ctx, postgresCollector.New(postgresConfig))
		require.NoError(t, err)

		_, err = normalDB.Exec("INSERT INTO test (name, ts) VALUES ('test1', now());")
		require.NoError(t, err)

		var (
			done     = make(chan struct{})
			commitCh = make(chan uint64)
		)

		time.AfterFunc(10*time.Second, func() {
			fmt.Println("canceling...")
			cancel()
		})

		go func() {
			for msg := range publishCh {
				eventData := string(msg.Get("WalData").([]byte))
				fmt.Println("event, data:", eventData)

				commitCh <- msg.Get("WalStart").(uint64)
			}
			close(commitCh)
		}()

		go func() {
			for offset := range commitCh {
				fmt.Println("--> commit, offset", pgx.FormatLSN(offset))
				offsetCh <- offset
				fmt.Println("<-- commit, offset", pgx.FormatLSN(offset))
			}

			close(offsetCh)
			close(done)
		}()

		<-done
	})
}
