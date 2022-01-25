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

package util

import "time"

// ErrFrequencyInterface err appear frequency interface
type ErrFrequencyInterface interface {
	// IsErrAlwaysAppear is error always appear
	IsErrAlwaysAppear(err error) bool

	// Release release error
	Release()
}

type errFrequency struct {
	err     error
	endTime int64
}

// NewErrFrequency new ErrorFrequency struct
func NewErrFrequency(err error) ErrFrequencyInterface {
	return &errFrequency{
		err:     err,
		endTime: time.Now().Add(10 * time.Minute).Unix(),
	}
}

// IsErrAlwaysAppear is error always appear
func (e *errFrequency) IsErrAlwaysAppear(err error) bool {
	if err == nil {
		return false
	}

	if e.err != nil && e.err.Error() == err.Error() {
		if time.Now().Unix() >= e.endTime {
			return true
		}
		return false
	}

	e.err = err
	e.endTime = time.Now().Add(10 * time.Minute).Unix()
	return false
}

// Release release error
func (e *errFrequency) Release() {
	e.err = nil
}
