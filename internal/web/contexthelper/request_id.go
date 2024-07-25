package contexthelper

import (
	"context"

	"github.com/ernestngugi/medvice-backend/internal/entities"
	"github.com/google/uuid"
)

func RequestId(ctx context.Context) string {
	existing := ctx.Value(entities.ContextKeyRequestID)
	if existing == nil {
		return ""
	}

	if val, ok := existing.(string); ok {
		u, err := uuid.Parse(val)
		if err != nil {
			return val
		}
		return u.String()
	}

	return ""
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, entities.ContextKeyRequestID, requestId)
}
