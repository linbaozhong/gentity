package log

import (
	"context"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
	"github.com/linbaozhong/gentity/pkg/app"
	"os"
	"path/filepath"
	"time"
)

var (
	level slog.Level
)

// Register slog logger
func Register(production bool) {
	_ = os.Mkdir("logs", os.ModePerm)

	slog.Configure(func(l *slog.SugaredLogger) {
		f := l.Formatter.(*slog.TextFormatter)
		f.SetTemplate("{{datetime}} [{{level}}] [{{caller}}] {{message}}\n")
		f.EnableColor = true
		f.TimeFormat = time.DateTime
		l.CallerSkip += 1

		errorHandler := handler.MustFileHandler(filepath.Join(".", "logs", "error.log"),
			handler.WithBuffMode(handler.BuffModeLine),
			handler.WithRotateTime(rotatefile.EveryDay),
			handler.WithLogLevels(slog.DangerLevels),
		)
		l.PushHandler(errorHandler)

		infoHandler := handler.MustFileHandler(filepath.Join(".", "logs", "info.log"),
			handler.WithBuffMode(handler.BuffModeLine),
			handler.WithRotateTime(rotatefile.EveryDay),
			handler.WithLogLevels(slog.NormalLevels),
		)
		l.PushHandler(infoHandler)
	})
	if production {
		level = slog.ErrorLevel
	} else {
		level = slog.TraceLevel
	}
	slog.SetLogLevel(level)
	// 注册关闭器
	app.RegisterServiceCloser(slog.Std())
}

// RegisterLogger slog logger
func RegisterLogger(ctx context.Context, production bool) {
	Register(production)
}

func Trace(args ...any) {
	slog.Trace(args...)
}
func Tracef(format string, args ...any) {
	slog.Tracef(format, args...)
}
func Debug(args ...any) {
	slog.Debug(args...)
}
func Debugf(format string, args ...any) {
	slog.Debugf(format, args...)
}
func Info(args ...any) {
	slog.Info(args...)
}
func Infof(format string, args ...any) {
	slog.Infof(format, args...)
}
func Notice(args ...any) {
	slog.Notice(args...)
}
func Noticef(format string, args ...any) {
	slog.Noticef(format, args...)
}
func Warn(args ...any) {
	slog.Warn(args...)
}
func Warnf(format string, args ...any) {
	slog.Warnf(format, args...)
}
func Error(args ...any) {
	slog.Error(args...)
}
func Errorf(format string, args ...any) {
	slog.Errorf(format, args...)
}

func Fatal(args ...any) {
	slog.Fatal(args...)
}
func Fatalf(format string, args ...any) {
	slog.Fatalf(format, args...)
}
func Panic(args ...any) {
	slog.Panic(args...)
}
func Panicf(format string, args ...any) {
	slog.Panicf(format, args...)
}

func setLevel(l slog.Level) {
	slog.SetLogLevel(l)
}
