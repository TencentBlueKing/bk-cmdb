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

package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

// Dir  exception file directory
var Dir = "/tmp"

// Exception use file record exception info
type Exception struct {
	flag      string
	fileName  string
	fileCount int
	directory string
	file      *os.File
	// taskID combin synchronize task flag and version
	taskID string
}

// NewException  new file exception struct
func NewException(synchronizieFlag string, version int64, fileCount int) (*Exception, error) {
	taskID := fmt.Sprintf("%s-%d", synchronizieFlag, version)
	file, err := os.Create(Dir + "/" + taskID)
	if err != nil {
		return nil, err
	}
	if fileCount <= 0 {
		fileCount = 10
	}

	return &Exception{
		directory: Dir,
		taskID:    taskID,
		file:      file,
		fileName:  taskID,
		fileCount: fileCount,
	}, nil
}

// WriteStringln write string, append \n
func (f *Exception) WriteStringln(str string) error {
	_, err := f.file.WriteString(str + "\n")
	if err != nil {
		blog.Errorf("Exception WriteStringln file %s error, err: %s", f.fileName, err.Error())
		return err
	}
	return nil
}

func (f *Exception) Write(ctx context.Context, exceptions []metadata.ExceptionResult) error {
	ln := []byte("\n")
	for _, exception := range exceptions {
		// don't case error
		line, _ := json.Marshal(exception)
		_, err := f.file.Write(line)
		if err != nil {
			blog.Errorf("Exception Write file %s error, err: %s, taskID:%s", f.fileName, err.Error(), f.taskID)
			return err
		}
		f.file.Write(ln)
	}
	return nil
}

// Destruct clean funcation
func (f *Exception) Destruct() {
	f.file.Close()
	go func() {
		files, err := ioutil.ReadDir(f.directory)
		if err != nil {
			blog.Errorf("Destruct remove log,excel ioutil.ReadDir error, error:%s,taskID:%s", err.Error(), f.taskID)
			return
		}
		var fileNames []string
		for _, file := range files {
			if strings.HasPrefix(file.Name(), f.flag) {
				fileNames =
					append(fileNames, file.Name())
			}
		}

		sort.Sort(sort.StringSlice(fileNames))
		lenFileNames := len(fileNames)
		if lenFileNames > f.fileCount {

		}
		deleleFileCnt := lenFileNames - f.fileCount
		for idx := 0; idx < deleleFileCnt; idx++ {
			os.Remove(fileNames[idx])
		}
	}()

}
