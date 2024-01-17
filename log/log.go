// Copyright The Notary Project Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package log provides logging functionality.
// 3rd party loggers that implement log.Logger: github.com/uber-go/zap.SugaredLogger
// and github.com/sirupsen/logrus.Logger.
package log

// Logger is implemented by users and/or 3rd party loggers.
// For example, github.com/uber-go/zap.SugaredLogger
// and github.com/sirupsen/logrus.Logger.
type Logger interface {
	// Debug logs a debug level message.
	Debug(args ...interface{})

	// Debugf logs a debug level message with format.
	Debugf(format string, args ...interface{})

	// Debugln logs a debug level message. Spaces are always added between
	// operands.
	Debugln(args ...interface{})

	// Info logs an info level message.
	Info(args ...interface{})

	// Infof logs an info level message with format.
	Infof(format string, args ...interface{})

	// Infoln logs an info level message. Spaces are always added between
	// operands.
	Infoln(args ...interface{})

	// Warn logs a warn level message.
	Warn(args ...interface{})

	// Warnf logs a warn level message with format.
	Warnf(format string, args ...interface{})

	// Warnln logs a warn level message. Spaces are always added between
	// operands.
	Warnln(args ...interface{})

	// Error logs an error level message.
	Error(args ...interface{})

	// Errorf logs an error level message with format.
	Errorf(format string, args ...interface{})

	// Errorln logs an error level message. Spaces are always added between
	// operands.
	Errorln(args ...interface{})
}
