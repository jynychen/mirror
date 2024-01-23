package logger

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Logger interface {
	io.Writer
	Debug(msg interface{}, keyvals ...interface{})
	Info(msg interface{}, keyvals ...interface{})
	Warn(msg interface{}, keyvals ...interface{})
	Error(msg interface{}, keyvals ...interface{})
	Fatal(msg interface{}, keyvals ...interface{})
	Print(msg interface{}, keyvals ...interface{})
}
type LoggerWthWriter struct {
	io.Writer
	*log.Logger
	lock sync.Mutex
}

func New(writer io.Writer) *LoggerWthWriter {
	logger := log.NewWithOptions(writer, log.Options{
		ReportTimestamp: true,
	})
	return &LoggerWthWriter{
		writer,
		logger,
		sync.Mutex{},
	}
}

func (l *LoggerWthWriter) Write(p []byte) (n int, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	buf := []byte{}
	for _, b := range bytes.Split(p, []byte("\r")) {
		buf = append(buf, []byte(time.Now().Format("2006-01-02 15:04:05 "))...)
		buf = append(buf, b...)
		buf = append(buf, '\r')
	}

	return l.Writer.Write(buf)
}
