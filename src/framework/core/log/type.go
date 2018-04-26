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

// Loger is the interface that must be implemented by every log
type Loger interface {

	// Info logs to the INFO logs.
	Info(args ...interface{})
	Infof(format string, args ...interface{})

	// Warning logs to the WARNING and INFO logs.
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	// Error logs to the ERROR、 WARNING and INFO logs.
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs.
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
