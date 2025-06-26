package dsl

import (
	"fmt"

	"github.com/rs/zerolog"
)

type DSLLogger interface {
	Info(msg string)
	Warn(msg string, err error)
	Error(msg string, err error)
}

func NewFMTLogger() DSLLogger {
	return &FMTLogger{}
}

type FMTLogger struct{}

func (l *FMTLogger) Info(msg string) {
	fmt.Println(msg)
}

func (l *FMTLogger) Warn(msg string, err error) {
	fmt.Println(msg, err)
}

func (l *FMTLogger) Error(msg string, err error) {
	fmt.Println(msg, err)
}

type ZeroLogger struct {
	logger *zerolog.Logger
}

func NewZeroLogger(logger *zerolog.Logger) DSLLogger {
	return &ZeroLogger{logger: logger}
}

func (l *ZeroLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *ZeroLogger) Warn(msg string, err error) {
	l.logger.Warn().Err(err).Msg(msg)
}

func (l *ZeroLogger) Error(msg string, err error) {
	l.logger.Error().Err(err).Msg(msg)
}
