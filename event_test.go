package puzzles

import (
	"testing"

	"puzzles/solutions/s01"
	"puzzles/solutions/s02"
	"puzzles/solutions/s03"
	"puzzles/solutions/s04"
	"puzzles/solutions/s05"

	"github.com/liennie/code-and-chill/pkg/eventtest"
)

func TestEvent(t *testing.T) {
	eventtest.Test(t, "event.yaml",
		s01.Solution,
		s02.Solution,
		s03.Solution,
		s04.Solution,
		s05.Solution,
	)
}
