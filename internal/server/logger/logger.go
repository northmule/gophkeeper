package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

type LogEntry struct {
	*zap.SugaredLogger
}

func NewLogger(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	appLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{appLogger.Sugar()}, nil
}

func (l *Logger) Print(v ...interface{}) {
	l.Info(v...)
}
func (l *Logger) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &LogEntry{
		l.SugaredLogger,
	}
}

func (l *LogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Infof("Request Information: Status: %d. Byte: %d. Headings: %#v. Time: %d. Additionally: %#v", status, bytes, header, elapsed, extra)
}
func (l *LogEntry) Panic(v interface{}, stack []byte) {
	l.Infof("Panic: %#v. Trace: %s", v, string(stack))
}
