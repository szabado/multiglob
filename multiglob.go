package multiglob

import (
	"github.com/szabado/multiglob/internal/parser"
)

type Builder struct {
	patterns map[string]*parser.Node
}

func New() *Builder {
	return &Builder{
		patterns: make(map[string]*parser.Node),
	}
}

func (m *Builder) AddPattern(name, pattern string) {
	m.patterns[name] = parser.Parse(pattern)
}

type MultiGlob struct{}
