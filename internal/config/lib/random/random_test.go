package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"short", 8},
		{"medium", 16},
		{"long", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.size)
			assert.Equal(t, tt.size, len(got))
		})
	}
}
