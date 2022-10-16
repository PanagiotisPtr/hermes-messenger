package loggingutils

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggerWithRequestID(
	ctx context.Context,
	l *zap.Logger,
) *zap.Logger {
	requestId, ok := ctx.Value("request-id").(uuid.UUID)
	if !ok {
		return l.With(
			zap.String("request-id", "none"),
		)
	}

	return l.With(
		zap.String("request-id", requestId.String()),
	)
}
