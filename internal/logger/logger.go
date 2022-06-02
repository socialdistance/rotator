package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"net/http"
	"rotator/internal/config"
	"time"
)

type Logger struct {
	zap *zap.Logger
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.zap.Debug(message, fields...)
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	l.zap.Info(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	l.zap.Error(message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.zap.Fatal(message, fields...)
}

func (l *Logger) LogHTTP(r *http.Request, code, length int) {
	l.zap.Info("HTTP logger:", zap.Strings(
		"HTTP", []string{r.RemoteAddr,
			time.Now().Format("01/Jan/2003:10:10:10 MST"),
			r.RequestURI,
			r.Proto,
			r.UserAgent()}),
		zap.Int("code", code),
		zap.Int("length", length))
}

func NewLogger(configLogger config.LoggerConf) (*Logger, error) {
	var cfg zap.Config

	marshal, err := json.Marshal(configLogger)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(marshal, &cfg); err != nil {
		log.Fatalf("Can't unmarshall logger %s", err)
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Can't build logger %s", err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")

	return &Logger{
		zap: logger,
	}, nil
}
