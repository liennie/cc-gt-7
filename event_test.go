package puzzles

import (
	"testing"

	"puzzles/solutions/puzzle01"
	"puzzles/solutions/puzzle02"
	"puzzles/solutions/puzzle03"
	"puzzles/solutions/puzzle04"
	"puzzles/solutions/puzzle05"

	"github.com/liennie/code-and-chill/pkg/eventtest"
)

func TestEvent(t *testing.T) {
	eventtest.Test(t, "event.yaml",
		puzzle01.Solution,
		puzzle02.Solution,
		puzzle03.Solution,
		puzzle04.Solution,
		puzzle05.Solution,
	)
}
