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
	requestId, ok := ctx.Value("requestId").(uuid.UUID)
	if !ok {
		return l.With(
			zap.String("requestId", "none"),
		)
	}

	return l.With(
		zap.String("requestId", requestId.String()),
	)
}
