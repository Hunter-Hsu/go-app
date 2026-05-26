package common

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

func NewContextWithTraceID(ctx context.Context, traceId string) context.Context {
	if traceId == "" {
		traceId = uuid.New().String()
	}
	return context.WithValue(ctx, traceIDKey, traceId)
}

func GetTraceID(ctx context.Context) string {
	traceId := ctx.Value(traceIDKey)
	if traceId == nil {
		return uuid.New().String()
	}
	return traceId.(string)
}
