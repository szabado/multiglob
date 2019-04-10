package multiglob

import (
	"testing"

	r "github.com/stretchr/testify/require"

	"github.com/szabado/multiglob/internal/parser"
)

func TestAddPattern(t *testing.T) {
	require := r.New(t)

	b := New()
	b.AddPattern("test", "pattern")

	require.Equal(1, len(b.patterns))
	require.Equal(parser.Parse("pattern"), b.patterns["test"])
}
