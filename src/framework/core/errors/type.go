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
 
package errors

import (
	"errors"
)

type ErrorsInterface interface {
	New() func(message string) error
}

type pkgError struct{}

func (pkgError) New() func(message string) error {
	return errors.New
}

// ErrNotSuppportedFunctionality returns an error cause the functionality is not supported
var ErrNotSuppportedFunctionality = errors.New("not supported functionality")

// ErrNotImplementedFunctionality returns an error cause the functionality is not implemented
var ErrNotImplementedFunctionality = errors.New("not implemented functionality")

// ErrDuplicateDataExisted returns an error cause the functionality is not supported
var ErrDuplicateDataExisted = errors.New("duplicated data existed")
