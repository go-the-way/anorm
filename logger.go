// Copyright 2022 anorm Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anorm

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	logger struct {
		mu         *sync.Mutex
		lo         *log.Logger
		low        io.Writer
		level      logLevel
		timeLayout string
	}
	logField struct{ k, v any }
	logLevel uint8
)

const (
	// LogLevelInfo => log info level
	LogLevelInfo = logLevel(1 + iota)
	// LogLevelDebug => log debug level
	LogLevelDebug
	// LogLevelError => log error level
	LogLevelError
	// LogLevelFatal => log fatal level
	LogLevelFatal
)

// Logger the global logger
var (
	Logger = &logger{
		mu:         &sync.Mutex{},
		lo:         log.New(os.Stdout, "[anorm] ", 0),
		level:      LogLevelError,
		timeLayout: time.RFC3339,
	}
	logLevelMap = map[logLevel]string{
		LogLevelInfo:  "INFO",
		LogLevelDebug: "DEBUG",
		LogLevelError: "ERROR",
		LogLevelFatal: "FATAL",
	}
	queryLog = func(name, sql string, ps []any) {
		Logger.Debug([]*logField{LogField("Name", name), LogField("SQL", sql), LogField("Parameter", ps)}, "")
	}
	queryErrorLog = func(name, sql string, ps []any, err error) {
		if err != nil {
			Logger.Error([]*logField{LogField("Name", name), LogField("SQL", sql), LogField("Parameter", ps)}, "err: %v", err)
		}
	}
	LogField = func(k, v any) *logField { return &logField{k, v} }
)

func (l *logger) getValidLevel(level logLevel) logLevel {
	if level < LogLevelInfo {
		level = LogLevelInfo
	}
	if level > LogLevelFatal {
		level = LogLevelFatal
	}
	return level
}

func (l *logger) getLogName(level logLevel) string {
	level = l.getValidLevel(level)
	return logLevelMap[level]
}

// SetLogLevel set global log level
func (l *logger) SetLogLevel(level logLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = l.getValidLevel(level)
}

// SetTimeLayout set time layout
func (l *logger) SetTimeLayout(layout string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeLayout = layout
}

// SetOutput set log output
func (l *logger) SetOutput(low io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lo.SetOutput(low)
}

// Info display info log
func (l *logger) Info(fields []*logField, message string, args ...any) {
	l.log(LogLevelInfo, fields, message, args...)
}

// Debug display info log
func (l *logger) Debug(fields []*logField, message string, args ...any) {
	l.log(LogLevelDebug, fields, message, args...)
}

// Error display info log
func (l *logger) Error(fields []*logField, message string, args ...any) {
	l.log(LogLevelError, fields, message, args...)
}

// Fatal display info log
func (l *logger) Fatal(fields []*logField, message string, args ...any) {
	l.log(LogLevelFatal, fields, message, args...)
}

func (l *logger) log(level logLevel, fields []*logField, message string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	level = l.getValidLevel(level)
	if level >= l.level {
		logArr := make([]string, 0)
		logArr = append(logArr, fmt.Sprintf("[%s]", time.Now().Format(l.timeLayout)))
		logArr = append(logArr, fmt.Sprintf("%s -", l.getLogName(level)))
		if fields != nil {
			for _, f := range fields {
				k, v := f.k, f.v
				logArr = append(logArr, fmt.Sprintf("%v{%v}", k, v))
			}
		}
		logArr = append(logArr, fmt.Sprintf(message, args...))
		logStr := strings.Join(logArr, " ")
		l.lo.Println(logStr)
	}
}
