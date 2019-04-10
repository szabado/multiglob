package multiglob

import (
	"sync"

	"github.com/szabado/multiglob/internal/parser"
)

type Builder struct {
	lock     *sync.Mutex
	patterns map[string]*parser.Node
}

func New() *Builder {
	return &Builder{
		patterns: make(map[string]*parser.Node),
		lock:     &sync.Mutex{},
	}
}

func (m *Builder) AddPattern(name, pattern string) {
	m.lock.Lock()
	m.patterns[name] = parser.Parse(pattern)
	m.lock.Unlock()
}
