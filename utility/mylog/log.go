// log
package mylog

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
	level  int
}

var (
	ErrUnknownLevel error = errors.New("unknown level")

	levelDebugName []string = []string{"FATAL", "ERROR", "WARNING", "INFO", "DEBUG"}
)

func NewLogger(prefix string, level int) *Logger {
	logger := &Logger{
		logger: log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags|log.Lshortfile),
		level:  level,
	}
	return logger
}

func (c *Logger) Debug(argvs ...interface{}) {
	c.outputln(4, argvs...)
}

func (c *Logger) Info(argvs ...interface{}) {
	c.outputln(3, argvs...)
}

func (c *Logger) Warning(argvs ...interface{}) {
	c.outputln(2, argvs...)
}

func (c *Logger) Error(argvs ...interface{}) {
	c.outputln(1, argvs...)
}

func (c *Logger) Fatalln(argvs ...interface{}) {
	c.outputln(0, argvs...)
	os.Exit(1)
}

func (c *Logger) outputln(level int, argvs ...interface{}) {
	if level < 0 || level > len(levelDebugName)-1 {
		log.Fatalln(ErrUnknownLevel)
	}
	if level > c.level {
		return
	}
	wb := bytes.NewBuffer(make([]byte, 0, 256))
	for _, v := range argvs {
		wb.WriteString(fmt.Sprintf("%v ", v))
	}
	wb.Truncate(wb.Len() - 1)
	ss := fmt.Sprintf("[%s] %s\n", levelDebugName[level], wb.String())
	c.logger.Output(3, ss)
}
