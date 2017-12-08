package file_test

import (
	"testing"

	"github.com/pagarme/warp-pipe/lib/tester/file"
)

func TestLoad(t *testing.T) {

	text := file.Load(t, "loremipson")
	file.MustMatch(t, "loremipson", text)
}
