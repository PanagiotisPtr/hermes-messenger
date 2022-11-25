package loggingutils

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetLoggerWith(
	ctx context.Context,
	l *zap.Logger,
	args ...func(context.Context, *zap.Logger) *zap.Logger,
) *zap.Logger {
	for _, f := range args {
		l = f(ctx, l)
	}

	return l
}

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

func LoggerWithUserID(
	ctx context.Context,
	l *zap.Logger,
) *zap.Logger {
	userId, ok := ctx.Value("user-id").(uuid.UUID)
	if !ok {
		return l.With(
			zap.String("user-id", "none"),
		)
	}

	return l.With(
		zap.String("user-id", userId.String()),
	)
}

// Provides the ZAP logger
func ProvideProductionLogger() *zap.Logger {
	logger, _ := zap.NewProduction()

	return logger
}
