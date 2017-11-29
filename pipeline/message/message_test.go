package message_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/pipeline/message"
)

func TestMessage(t *testing.T) {

	testTable := []struct {
		name  string
		use   bool
		key   string
		value string
	}{
		{"Valid", true, "testKey", "testValue"},
		{"Invalid", false, "nilKey", ""},
	}

	payload := message.Payload{}
	for _, testData := range testTable {
		if !testData.use {
			continue
		}
		payload[testData.key] = testData.value
	}

	m := message.New(payload)
	require.NotNil(t, m)
	require.True(t, m.Timestamp().Before(time.Now()))

	for _, testData := range testTable {
		t.Run(testData.name, func(t *testing.T) {
			if testData.use {
				require.Equal(t, testData.value, m.Get(testData.key), "test failed:", testData.name)
			} else {
				require.Nil(t, m.Get(testData.key), "test failed:", testData.name)
			}
		})
	}
}
