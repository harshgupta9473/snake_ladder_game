package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
type ZapLogger interface{
	LogInfo(msg string, fields ...zap.Field)
	Close()
	Start()
}
func NewLogger(field string, name string) (ZapLogger) {
	logDir := `C:\Coding\WebD\Bebetta\logs`
	os.MkdirAll(logDir, 0755)

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", name, timestamp))

	file, err := os.Create(logFilePath)
	if err != nil {
		return nil
	}

	writer := zapcore.AddSync(file)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "msg"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writer,
		zap.InfoLevel,
	)

	zapLogger := zap.New(core).With(zap.String(field, name))

	return &Logger{
		Field:    field,
		Logger:   zapLogger,
		interval: 1*time.Hour,
		Name: name,
	}
}

type Logger struct {
	Field string
	Name string
	Logger      *zap.Logger
	file        *os.File
	mutex       sync.Mutex
	interval    time.Duration
}

func (l *Logger) LogInfo(msg string, fields ...zap.Field) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.Logger.Info(msg, fields...)
}

func (l *Logger)Start(){
	go l.doRotateInInterval()
}



func (l *Logger) doRotateInInterval() {
	for {
		time.Sleep(l.interval)

		l.mutex.Lock()

		// Close old file if open
		if l.file != nil {
			_ = l.file.Close()
		}

		// Create new log file with timestamp
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logDir := `C:\Coding\WebD\Bebetta\logs`
		logFilePath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", l.Name, timestamp))

		file, err := os.Create(logFilePath)
		if err != nil {
			fmt.Printf("Log rotation failed: %v\n", err)
			l.mutex.Unlock()
			continue
		}

		writer := zapcore.AddSync(file)
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.LevelKey = "level"
		encoderCfg.MessageKey = "msg"

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			writer,
			zap.InfoLevel,
		)

		// Create new logger with service field
		newLogger := zap.New(core).With(zap.String(l.Field, l.Name))

		// Replace old logger and file
		l.Logger = newLogger
		l.file = file

		l.mutex.Unlock()
	}
}


func (l *Logger) Close() {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    if l.file != nil {
        _ = l.file.Close()
    }
    _ = l.Logger.Sync() // flush zap
}