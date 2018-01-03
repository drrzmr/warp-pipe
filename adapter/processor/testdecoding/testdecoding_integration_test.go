package testdecoding_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pagarme/warp-pipe/adapter/collector/postgres"
	"github.com/pagarme/warp-pipe/adapter/processor/testdecoding"
	"github.com/pagarme/warp-pipe/lib/log"
	parser "github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate"
	tester "github.com/pagarme/warp-pipe/lib/tester/postgres"
	"github.com/pagarme/warp-pipe/pipeline"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Setup(log.Test)
}

func TestIntegrationPostgresCollectorTestDecodingProcessor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}

	dockerConfig := replicate.CreateTestDockerConfig(t)
	postgresConfig := replicate.CreateTestPostgresConfig(t)

	normal, deferFn := tester.DockerRun(t, dockerConfig, &postgresConfig)
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

		var (
			run         = pipeline.NewRunner("test")
			ctx, cancel = context.WithCancel(context.Background())
			collector   = postgres.New(postgresConfig)
			processor   = testdecoding.New()
		)

		outOfCollector, offsetCh, err := run.Collector(ctx, collector)
		require.NoError(t, err)
		close(offsetCh)

		outOfProcessor, err := run.Processor(ctx, processor, outOfCollector)
		require.NoError(t, err)

		_, err = normalDB.Exec("INSERT INTO test (name, ts) VALUES ('test1', now());")
		require.NoError(t, err)

		time.AfterFunc(10*time.Second, func() {
			fmt.Println("[canceling...]")
			cancel()
		})

		for msg := range outOfProcessor {

			transaction := msg.Get("transaction").(parser.Transaction)
			jsonTransaction, err := json.MarshalIndent(transaction, "", "\t")
			require.NoError(t, err)

			fmt.Println(string(jsonTransaction))
		}
	})
}
