package logger

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Logger interface {
	io.Writer
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}
type LoggerWthWriter struct {
	log.Logger
	lock sync.Mutex
}

func New() *LoggerWthWriter {
	return &LoggerWthWriter{
		*log.New(os.Stdout, "", log.LstdFlags),
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

	return l.Logger.Writer().Write(buf)
}
