package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

//TODO пока никак не юзается, нужно прикрутить к деливери слою, должен быть реквестid
type Logger interface {
	//GetRequestIdKey() int
	//GetRequestIdFromContext(ctx context.Context) string
	StartRequest(r http.Request, requestId string)
	EndRequest(start time.Time, ctx context.Context)
	HttpInfo(ctx context.Context, msg string, status int)
	HttpLogWarning(ctx context.Context, pkg string, funcName string, warn string)
	HttpLogError(ctx context.Context, pkg string, funcName string, err error)
}

const requestIdKey int = 1

type SimpleLogger struct {
	*logrus.Logger
}

func NewSimpleLogger(writer io.Writer, formatter logrus.Formatter) SimpleLogger {
	baseLogger := logrus.New()
	simpleLogger := SimpleLogger{baseLogger}
	simpleLogger.SetFormatter(formatter)
	simpleLogger.SetOutput(writer)
	return simpleLogger
}

func NewJsonFormatSimpleLogger(writer io.Writer) SimpleLogger {
	formatter := new(logrus.JSONFormatter)
	formatter.TimestampFormat = "03-05-2012 13:56:33"
	return NewSimpleLogger(writer, formatter)
}

func NewTextFormatSimpleLogger(writer io.Writer) SimpleLogger {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "03-05-2012 13:56:33"
	return NewSimpleLogger(writer, formatter)
}

func (logger SimpleLogger) GetRequestIdKey() int {
	return requestIdKey
}

func (logger SimpleLogger) GetRequestIdFromContext(ctx context.Context) string {
	requestId, ok := ctx.Value(logger.GetRequestIdKey()).(string)
	if !ok {
		logger.WithFields(logrus.Fields{
			"id":       "NO_ID",
			"package":  "logger",
			"function": "GetRequestIdFromContext",
		}).Warn("can't get request id from context")
		return ""
	}
	return requestId
}

func (logger SimpleLogger) StartRequest(r http.Request, requestId string) {
	logger.WithFields(logrus.Fields{
		"id":         requestId,
		"usr_addr":   r.RemoteAddr,
		"req_URI":    r.RequestURI,
		"method":     r.Method,
		"user_agent": r.UserAgent(),
	}).Info("request started")
}

func (logger SimpleLogger) EndRequest(start time.Time, ctx context.Context) {
	logger.WithFields(logrus.Fields{
		"id":              logger.GetRequestIdFromContext(ctx),
		"elapsed_time,μs": time.Since(start).Microseconds(),
	}).Info("request ended")
}

func (logger SimpleLogger) HttpInfo(ctx context.Context, msg string, status int) {
	logger.WithFields(logrus.Fields{
		"id":     logger.GetRequestIdFromContext(ctx),
		"status": status,
	}).Info(msg)
}

func (logger SimpleLogger) HttpLogWarning(ctx context.Context, pkg string, funcName string, warn string) {
	logger.WithFields(logrus.Fields{
		"id":       logger.GetRequestIdFromContext(ctx),
		"package":  pkg,
		"function": funcName,
	}).Warn(warn)
}

func (logger SimpleLogger) HttpLogError(ctx context.Context, pkg string, funcName string, err error) {
	logger.WithFields(logrus.Fields{
		"id":       logger.GetRequestIdFromContext(ctx),
		"package":  pkg,
		"function": funcName,
	}).Error(err)
}
