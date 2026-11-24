package store

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMemoryGetSubscriptionEdgeCases covers the lookups that are not the happy
// path: the memory store must answer ErrNotFound for every secret it never
// issued, without returning a subscription or panicking.
func TestMemoryGetSubscriptionEdgeCases(t *testing.T) {
	t.Parallel()

	m := NewMemory()
	cases := []struct {
		name   string
		secret string
	}{
		{name: "empty secret", secret: ""},
		{name: "whitespace secret", secret: "   "},
		{name: "very long secret", secret: strings.Repeat("x", 4096)},
		{name: "issued-looking but unknown secret", secret: "soro_2f8c1d9e4b7a6c5d"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sub, err := m.GetSubscription(context.Background(), tc.secret)
			require.ErrorIs(t, err, ErrNotFound)
			assert.Nil(t, sub)
		})
	}

	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		sub, err := m.GetSubscription(ctx, "any")
		require.ErrorIs(t, err, context.Canceled)
		assert.Nil(t, sub)
	})
}
