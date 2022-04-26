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
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestLogger_getValidLevel(t *testing.T) {
	for _, tc := range []struct {
		testName string
		set      logLevel
		expect   logLevel
	}{
		{"Test-0", logLevel(0), LogLevelInfo},
		{"Test-1", logLevel(1), LogLevelInfo},
		{"Test-2", logLevel(2), LogLevelDebug},
		{"Test-3", logLevel(3), LogLevelError},
		{"Test-4", logLevel(4), LogLevelFatal},
		{"Test-5", logLevel(5), LogLevelFatal},
		{"Test-6", logLevel(6), LogLevelFatal},
		{"Test-100", logLevel(100), LogLevelFatal},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			if expect := Logger.getValidLevel(tc.set); expect != tc.expect {
				t.Error("test failed")
			}
		})
	}
}

func TestLogger_SetLogLevel(t *testing.T) {
	for _, tc := range []struct {
		testName string
		set      logLevel
		expect   logLevel
	}{
		{"Test-0", logLevel(0), LogLevelInfo},
		{"Test-1", logLevel(1), LogLevelInfo},
		{"Test-2", logLevel(2), LogLevelDebug},
		{"Test-3", logLevel(3), LogLevelError},
		{"Test-4", logLevel(4), LogLevelFatal},
		{"Test-5", logLevel(5), LogLevelFatal},
		{"Test-6", logLevel(6), LogLevelFatal},
		{"Test-100", logLevel(100), LogLevelFatal},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			Logger.SetLogLevel(tc.set)
			if Logger.level != tc.expect {
				t.Error("test failed")
			}
		})
	}
}

func TestLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	layout := "2006-01-02"
	Logger.SetTimeLayout(layout)
	Logger.SetLogLevel(LogLevelInfo)
	Logger.SetOutput(buf)
	for _, tc := range []struct {
		testName string
		logLevel logLevel
		expect   string
		fields   []*logField
		message  string
		args     any
	}{
		{"Test-info-no-fields-message-no-args", LogLevelInfo, fmt.Sprintf("[anorm] [%s] INFO - hello world\n", time.Now().Format(layout)), nil, "hello world", nil},
		{"Test-info-no-fields-message-args", LogLevelInfo, fmt.Sprintf("[anorm] [%s] INFO - hello world-10\n", time.Now().Format(layout)), nil, "hello world-%d", 10},
		{"Test-info-fields-message-no-args", LogLevelInfo, fmt.Sprintf("[anorm] [%s] INFO - ID{1000} NAME{Apple} hello world\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world", nil},
		{"Test-info-fields-message-args", LogLevelInfo, fmt.Sprintf("[anorm] [%s] INFO - ID{1000} NAME{Apple} hello world-10\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world-%d", 10},

		{"Test-debug-no-fields-message-no-args", LogLevelDebug, fmt.Sprintf("[anorm] [%s] DEBUG - hello world\n", time.Now().Format(layout)), nil, "hello world", nil},
		{"Test-debug-no-fields-message-args", LogLevelDebug, fmt.Sprintf("[anorm] [%s] DEBUG - hello world-10\n", time.Now().Format(layout)), nil, "hello world-%d", 10},
		{"Test-debug-fields-message-no-args", LogLevelDebug, fmt.Sprintf("[anorm] [%s] DEBUG - ID{1000} NAME{Apple} hello world\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world", nil},
		{"Test-debug-fields-message-args", LogLevelDebug, fmt.Sprintf("[anorm] [%s] DEBUG - ID{1000} NAME{Apple} hello world-10\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world-%d", 10},

		{"Test-error-no-fields-message-no-args", LogLevelError, fmt.Sprintf("[anorm] [%s] ERROR - hello world\n", time.Now().Format(layout)), nil, "hello world", nil},
		{"Test-error-no-fields-message-args", LogLevelError, fmt.Sprintf("[anorm] [%s] ERROR - hello world-10\n", time.Now().Format(layout)), nil, "hello world-%d", 10},
		{"Test-error-fields-message-no-args", LogLevelError, fmt.Sprintf("[anorm] [%s] ERROR - ID{1000} NAME{Apple} hello world\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world", nil},
		{"Test-error-fields-message-args", LogLevelError, fmt.Sprintf("[anorm] [%s] ERROR - ID{1000} NAME{Apple} hello world-10\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world-%d", 10},

		{"Test-fatal-no-fields-message-no-args", LogLevelFatal, fmt.Sprintf("[anorm] [%s] FATAL - hello world\n", time.Now().Format(layout)), nil, "hello world", nil},
		{"Test-fatal-no-fields-message-args", LogLevelFatal, fmt.Sprintf("[anorm] [%s] FATAL - hello world-10\n", time.Now().Format(layout)), nil, "hello world-%d", 10},
		{"Test-fatal-fields-message-no-args", LogLevelFatal, fmt.Sprintf("[anorm] [%s] FATAL - ID{1000} NAME{Apple} hello world\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world", nil},
		{"Test-fatal-fields-message-args", LogLevelFatal, fmt.Sprintf("[anorm] [%s] FATAL - ID{1000} NAME{Apple} hello world-10\n", time.Now().Format(layout)), []*logField{LogField("ID", "1000"), LogField("NAME", "Apple")}, "hello world-%d", 10},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			if tc.args != nil {
				Logger.log(tc.logLevel, tc.fields, tc.message, tc.args)
			} else {
				Logger.log(tc.logLevel, tc.fields, tc.message)
			}
			if s := buf.String(); s != tc.expect {
				t.Error("test failed")
			}
			buf.Reset()
		})
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	layout := "2006-01-02"
	Logger.SetTimeLayout(layout)
	Logger.SetLogLevel(LogLevelInfo)
	Logger.SetOutput(buf)
	Logger.Info(nil, "hello world")
	expect := fmt.Sprintf("[anorm] [%s] INFO - hello world\n", time.Now().Format(layout))
	if s := buf.String(); s != expect {
		t.Error("test failed")
	}
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	layout := "2006-01-02"
	Logger.SetTimeLayout(layout)
	Logger.SetLogLevel(LogLevelInfo)
	Logger.SetOutput(buf)
	Logger.Debug(nil, "hello world")
	expect := fmt.Sprintf("[anorm] [%s] DEBUG - hello world\n", time.Now().Format(layout))
	if s := buf.String(); s != expect {
		t.Error("test failed")
	}
}

func TestLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	layout := "2006-01-02"
	Logger.SetTimeLayout(layout)
	Logger.SetLogLevel(LogLevelInfo)
	Logger.SetOutput(buf)
	Logger.Error(nil, "hello world")
	expect := fmt.Sprintf("[anorm] [%s] ERROR - hello world\n", time.Now().Format(layout))
	if s := buf.String(); s != expect {
		t.Error("test failed")
	}
}

func TestLogger_Fatal(t *testing.T) {
	buf := &bytes.Buffer{}
	layout := "2006-01-02"
	Logger.SetTimeLayout(layout)
	Logger.SetLogLevel(LogLevelInfo)
	Logger.SetOutput(buf)
	Logger.Fatal(nil, "hello world")
	expect := fmt.Sprintf("[anorm] [%s] FATAL - hello world\n", time.Now().Format(layout))
	if s := buf.String(); s != expect {
		t.Error("test failed")
	}
}
