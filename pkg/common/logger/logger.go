package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type ctxKey struct{}

var loggerKey = &ctxKey{}

const (
	usernameKey = "username"
	// requestIDKey = "request_id"
)

type Config struct {
	Level         string
	IsDevelopment bool
}

type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Err(err error) Logger
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	Sync() error
}

type zapLogger struct {
	logger *zap.Logger
}

// GetDefaultConfig is only used in case of an error getting the logger
// Using this config in production is an issue
func GetDefaultLogger() Logger {
	return NewLogger(Config{
		IsDevelopment: false,
		Level:         "info",
	})
}

func NewLogger(config Config) Logger {
	var zapConfig zap.Config

	if config.IsDevelopment {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Encoding = "json"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		zapConfig.DisableStacktrace = true
		zapConfig.DisableCaller = true
	}

	switch config.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		panic("invalid log level")
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{
		logger: logger,
	}
}

func SetLoggerToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func SetLoggerMiddlewareEcho(e *echo.Echo, baseLogger Logger) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := baseLogger
			ctx := c.Request().Context()
			// if requestID := tracing.GetTracingIDFromContext(ctx); requestID != "" {
			// 	logger = logger.WithField(requestIDKey, requestID)
			// }

			ctx = context.WithValue(ctx, loggerKey, logger)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
}

func SetUserDetailsToLogger(c echo.Context, username string) {
	ctx := c.Request().Context()
	logger := GetLogger(ctx)
	newLogger := logger.WithField(usernameKey, username)
	ctx = context.WithValue(ctx, loggerKey, newLogger)
	c.SetRequest(c.Request().WithContext(ctx))
}

func SetHTTPLoggerMiddlewareEcho(e *echo.Echo, urlSkipper func(c echo.Context) bool) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogLatency:  true,
		LogMethod:   true,
		HandleError: true,
		Skipper:     urlSkipper,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			ctx := c.Request().Context()
			log := GetLogger(ctx)
			log.WithFields(map[string]interface{}{
				"method":  v.Method,
				"uri":     v.URI,
				"status":  v.Status,
				"latency": v.Latency,
			}).Infof("HTTP Request")

			return nil
		},
	}))
}

func GetLogger(c context.Context) Logger {
	if l, ok := c.Value(loggerKey).(Logger); ok {
		return l
	}

	log.Error("logger not found in context")
	l := GetDefaultLogger()
	return l
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l *zapLogger) WithField(key string, value interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(zap.Any(key, value)),
	}
}

// Ensure a key can only be present once in the fields map
func (l *zapLogger) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return l
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &zapLogger{
		logger: l.logger.With(zapFields...),
	}
}

func (l *zapLogger) Err(err error) Logger {
	if err != nil {
		return l.WithField("error", err.Error())
	}

	return l
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
