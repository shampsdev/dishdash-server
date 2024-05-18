package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer_WithPort(t *testing.T) {
	s := NewServer(WithPort(1234))
	assert.Equal(t, ":1234", s.Addr)
}
