package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// Logger represents the logger configuration.
type Logger struct {
	log *logrus.Logger
}

var log = &Logger{}

// NewLogger creates a new logger instance with default settings.
func NewLogger() *Logger {
	return &Logger{log: logrus.New()}
}

// SetLogger initializes the logger with the desired configuration.
func (l *Logger) SetLogger() {
	// Log as JSON instead of the default ASCII formatter.
	l.log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout and a log file.
	l.log.SetOutput(&lumberjack.Logger{
		Filename:   "log/orderserv.log",
		MaxSize:    10 * 1024 * 1024,
		MaxBackups: 3,
		MaxAge:     2,
	})

	//l.log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	l.log.SetLevel(logrus.InfoLevel)

	// conn, err := net.Dial("tcp", "0.0.0.0:5044")
	// if err != nil {
	// 	l.log.Fatal(err)
	// }

	// hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "userserv"}))

	// l.log.Hooks.Add(hook)

	log = l
	// Log the startup message.
	log.Info("Logger initialized successfully")

}

// Log returns the configured logger.
func Log() *Logger {
	return log
}

// Info logs an informational message.
func (l *Logger) Info(message string) {
	l.log.Info(message)
}

// Warn logs a warning message.
func (l *Logger) Warn(message string) {
	l.log.Warn(message)
}

// Error logs an error message.
func (l *Logger) Error(message string) {
	l.log.Error(message)
}

// Fatal logs a fatal message and exits the program.
func (l *Logger) Fatal(message string) {
	l.log.Fatal(message)
	os.Exit(1)
}
