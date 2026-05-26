package common

import (
	"context"
	"testing"

	"pgregory.net/rapid"
)

func TestProperty_GetTraceID_RoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		id := rapid.StringMatching(`[a-zA-Z0-9_-]+`).Draw(t, "traceID")
		ctx := context.Background()
		ctx = NewContextWithTraceID(ctx, id)
		got := GetTraceID(ctx)
		if got != id {
			t.Fatalf("round-trip failed: NewContextWithTraceID then GetTraceID: got %q, want %q", got, id)
		}
	})
}

func TestProperty_GetTraceID_NoInterference(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		id := rapid.StringMatching(`[a-zA-Z0-9]+`).Draw(t, "traceID")
		type sentinelKey struct{}
		ctx := context.WithValue(context.Background(), sentinelKey{}, "other-value")
		ctx = NewContextWithTraceID(ctx, id)
		if got := ctx.Value(sentinelKey{}); got != "other-value" {
			t.Fatalf("trace ID key interfered with other context value: got %v", got)
		}
	})
}

func TestGetTraceID_EmptyContext_ReturnsNonEmptyString(t *testing.T) {
	ctx := context.Background()
	id := GetTraceID(ctx)
	if id == "" {
		t.Fatal("GetTraceID on empty context should return a non-empty fallback UUID")
	}
}
