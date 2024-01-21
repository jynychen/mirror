package logger

import (
	"bytes"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	log.Logger
	lock sync.Mutex
}

func New() *Logger {
	return &Logger{
		*log.New(os.Stdout, "", log.LstdFlags),
		sync.Mutex{},
	}
}

func (l *Logger) Write(p []byte) (n int, err error) {
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
