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
 
package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// AtomicFile behaves like os.File, but does an atomic rename operation at Close.
type AtomicFile struct {
	*os.File
	path string
}

// AtomicFileNew creates a new temporary file that will replace the file at the given
// path when Closed.
func AtomicFileNew(path string, mode os.FileMode) (*AtomicFile, error) {
	f, err := ioutil.TempFile(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(f.Name(), mode); err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, closeErr
		}

		if removeErr := os.Remove(f.Name()); removeErr != nil {
			return nil, removeErr
		}
		return nil, err
	}
	return &AtomicFile{File: f, path: path}, nil
}

// Close the file replacing the configured file.
func (f *AtomicFile) Close() error {
	if err := f.File.Close(); err != nil {
		if removeErr := os.Remove(f.File.Name()); removeErr != nil {
			return removeErr
		}
		return err
	}
	if err := os.Rename(f.Name(), f.path); err != nil {
		return err
	}
	return nil
}

// Abort closes the file and removes it instead of replacing the configured
// file. This is useful if after starting to write to the file you decide you
// don't want it anymore.
func (f *AtomicFile) Abort() error {
	if err := f.File.Close(); err != nil {
		if removeErr := os.Remove(f.Name()); removeErr != nil {
			return removeErr
		}
		return err
	}
	if err := os.Remove(f.Name()); err != nil {
		return err
	}
	return nil
}
