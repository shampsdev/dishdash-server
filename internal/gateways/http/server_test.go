package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer_WithPort(t *testing.T) {
	s := NewServer(WithPort(1234))
	assert.Equal(t, ":1234", s.Addr)
}
