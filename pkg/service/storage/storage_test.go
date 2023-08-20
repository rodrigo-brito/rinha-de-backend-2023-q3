package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_UUID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"test", "098f6bcd-4621-d373-cade-4e832627b4f6"},
		{"nome muito grande para ver se quebra", "e607f008-d17c-6675-e0c8-66133b2f77f3"},
		{" ", "7215ee9c-7d9d-c229-d292-1a40e899ec5f"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			s := NewStorage()
			require.Equal(t, tt.want, s.UUID(tt.input))
		})
	}
}
