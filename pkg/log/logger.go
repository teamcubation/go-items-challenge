package log

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const RequestIDKey = "x-request-id"

var baseLogger *logrus.Logger

type loggerKey struct{}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	baseLogger = logger

	return baseLogger
}

func getLogger() *logrus.Logger {
	if baseLogger == nil {
		baseLogger = NewLogger()
	}

	return baseLogger
}

func Context(r *http.Request) context.Context {
	reqID := r.Header.Get(RequestIDKey)
	if reqID == "" {
		reqID = uuid.New().String()
	}

	loggerWithRequestID := getLogger().WithField("request_id", reqID)

	return context.WithValue(r.Context(), loggerKey{}, loggerWithRequestID)
}

func GetFromContext(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(loggerKey{}).(*logrus.Entry)
	if !ok {
		return defaultLogger()
	}

	return logger
}

func Println(ctx context.Context, args ...interface{}) {
	GetFromContext(ctx).Println(args...)
}

func defaultLogger() *logrus.Entry {
	return getLogger().WithField("default", "true")
}
