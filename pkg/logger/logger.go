package logger

import "go.uber.org/zap"

type Interface interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Sync()
}

type Logger struct {
	logger *zap.SugaredLogger
}

func New() (*Logger, error) {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		return nil, err
	}
	return &Logger{logger.Sugar()}, nil
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Infow(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warnw(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Errorw(msg, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}){
	l.logger.Errorf(msg,args)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalw(msg, args...)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}
