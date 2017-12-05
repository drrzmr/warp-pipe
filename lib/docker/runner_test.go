package docker_test

import (
	"testing"

	"github.com/pagarme/warp-pipe/lib/docker"
	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {

	table := []struct {
		url      string
		image    string
		tag      string
		expected string
	}{
		{url: "warp-pipe", image: "postgres-server", tag: "9.5.6", expected: "warp-pipe/postgres-server:9.5.6"},
	}

	t.Run("Config/ImageName", func(t *testing.T) {

		for _, test := range table {

			t.Run(test.expected, func(t *testing.T) {
				runner := docker.NewRunner(docker.Config{
					URL:   test.url,
					Image: test.image,
					Tag:   test.tag,
				})
				require.NotNil(t, runner)
				config := runner.Config()
				require.NotNil(t, config)
				require.Equal(t, test.expected, config.ImageName())
			})
		}
	})

}
