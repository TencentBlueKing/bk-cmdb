/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

var (
	// Info logs to waring error fatal
	Info func(args ...interface{})

	// Infof logs to  warning error fatal
	Infof func(format string, args ...interface{})

	// Warning logs to warning error fatal.
	Warning func(args ...interface{})

	// Warningf logs to warning error fatal.
	Warningf func(format string, args ...interface{})

	// Error logs to the error fatal
	Error func(args ...interface{})

	// Errorf logs to the error fatal
	Errorf func(format string, args ...interface{})

	// Fatal logs to the fatal
	Fatal func(args ...interface{})

	// Fatalf logs to fatal
	Fatalf func(format string, args ...interface{})
)

// SetLoger set a new loger instance
func SetLoger(target *Logger) {

	// Infof print info
	Infof = target.Infof

	// Info print info
	Info = target.Info

	// Warning print warning
	Warning = target.Warning

	// Warningf print the warningf
	Warningf = target.Warningf

	// Error print the error info
	Error = target.Error

	// Errorf print the error info
	Errorf = target.Errorf

	// Fatal print the fatal
	Fatal = target.Fatal

	// Fatalf print the fatal
	Fatalf = target.Fatalf
}

// Logger implements the Loger interface
type Logger struct {
	// Info logs to the INFO logs.
	Info  func(args ...interface{})
	Infof func(format string, args ...interface{})

	// Warning logs to the WARNING and INFO logs.
	Warning  func(args ...interface{})
	Warningf func(format string, args ...interface{})

	// Error logs to the ERROR、 WARNING and INFO logs.
	Error  func(args ...interface{})
	Errorf func(format string, args ...interface{})

	// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs.
	Fatal  func(args ...interface{})
	Fatalf func(format string, args ...interface{})
}
