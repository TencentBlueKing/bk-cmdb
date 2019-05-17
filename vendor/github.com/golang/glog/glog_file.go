// Go support for leveled logs, analogous to https://code.google.com/p/google-glog/
//
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// File I/O for logs.

package glog

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// MaxSize is the maximum size of a log file in bytes.
var maxSizeFlag = flag.Uint64("log_max_size", 500, "Max size (MB) per file.")
var maxSize uint64 = 0

func MaxSize() uint64 {
	if maxSize == 0 {
		maxSize = *maxSizeFlag * 1024 * 1024
	}
	return maxSize
}

// MaxNum is the maximum of log files for one thread.
var maxNumFlag = flag.Int("log_max_num", 6, "Max num of file. The oldest will be removed if there is a extra file created.")

func MaxNum() int {
	return *maxNumFlag
}

// fileInfo contains log filename and its timestamp.
type fileInfo struct {
	name      string
	timestamp string
}

// fileInfoList implements Interface interface in sort. For
// sorting a list of fileInfo
type fileInfoList []fileInfo

func (b fileInfoList) Len() int           { return len(b) }
func (b fileInfoList) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b fileInfoList) Less(i, j int) bool { return b[i].timestamp < b[j].timestamp }

// fileBlock is a block of chain in logKeeper.
type fileBlock struct {
	fileInfo
	next *fileBlock
}

// logKeeper maintains a chain of each level log file. Its head
// is the earliest file while its tail is the oldest. It remains
// up to MaxNum() files, and the extra added will lead to delete
// the oldest. At first it load from logDir and take existing files
// into the chain. And remove the part over MaxNum().
type logKeeper struct {
	dir      string
	onceLoad sync.Once
	head     map[string]*fileBlock
	tail     map[string]*fileBlock
	total    map[string]int
}

func (lk *logKeeper) add(tag string, newBlock *fileBlock) (ok bool) {
	block, ok := lk.tail[tag]
	if !ok {
		return
	}
	if block == nil {
		lk.head[tag] = newBlock
	} else {
		if block.name == newBlock.name {
			return false
		}
		block.next = newBlock
	}
	lk.tail[tag] = newBlock
	lk.total[tag]++
	for lk.total[tag] > MaxNum() {
		lk.remove(tag)
	}
	return
}

func (lk *logKeeper) remove(tag string) (ok bool) {
	block, ok := lk.head[tag]
	if !ok || lk.total[tag] == 0 {
		return
	}
	lk.removeFile(block.name)
	lk.head[tag] = block.next
	block = nil // for GC
	lk.total[tag]--
	return
}

func (lk *logKeeper) removeFile(name string) error {
	return os.Remove(filepath.Join(lk.dir, name))
}

func (lk *logKeeper) load() {
	_dir, err := ioutil.ReadDir(lk.dir)
	if err != nil {
		return
	}

	reg := logNameReg()
	tmp := make(map[string]fileInfoList)
	for _, fi := range _dir {
		if fi.IsDir() {
			continue
		}

		result := reg.FindStringSubmatch(fi.Name())
		if result == nil {
			continue
		}

		name, tag, timestamp := result[0], result[1], result[2]
		if tmp[tag] == nil {
			tmp[tag] = make([]fileInfo, 0, len(_dir))
		}
		tmp[tag] = append(tmp[tag], fileInfo{name: name, timestamp: timestamp})
	}

	for tag, blockList := range tmp {
		sort.Sort(blockList)
		for i, block := range blockList {
			if i < MaxNum() {
				fb := &fileBlock{
					fileInfo: fileInfo{name: block.name, timestamp: block.timestamp},
					next:     nil,
				}
				if i == 0 {
					lk.head[tag] = fb
				} else {
					lk.tail[tag].next = fb
				}
				lk.tail[tag] = fb
				lk.total[tag]++
			} else {
				lk.removeFile(block.name)
			}
		}
	}
}

// logDirs lists the candidate directories for new log files.
var logDirs []*logKeeper

// If non-empty, overrides the choice of directory in which to write logs.
// See createLogDirs for the full list of possible destinations.
var logDir = flag.String("log_dir", "", "If non-empty, write log files in this directory")

func createLogDirs() {
	var dirs []string
	if *logDir != "" {
		dirs = append(dirs, *logDir)
	}
	dirs = append(dirs, os.TempDir())

	for _, dir := range dirs {
		head := make(map[string]*fileBlock)
		tail := make(map[string]*fileBlock)
		total := make(map[string]int)
		for _, name := range severityName {
			head[name] = nil
			tail[name] = nil
			total[name] = 0
		}

		logDirs = append(logDirs, &logKeeper{dir: dir, head: head, tail: tail, total: total})
	}
}

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func logName(tag string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s.%s.log.%s.%04d%02d%02d-%02d%02d%02d.%d",
		program,
		host,
		userName,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, program + "." + tag
}

// logNameReg returns a regexp object for match log file name.
func logNameReg() *regexp.Regexp {
	reg, _ := regexp.Compile(fmt.Sprintf(`^%s\..+\..+\.log.(%s)\.(\d{8}-\d{6})\.\d+$`,
		program,
		strings.Join(severityName, "|")))
	return reg
}

var onceLogDirs sync.Once

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	onceLogDirs.Do(createLogDirs)
	if len(logDirs) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}
	name, link := logName(tag, t)
	var lastErr error
	for _, lk := range logDirs {
		lk.onceLoad.Do(lk.load)
		fname := filepath.Join(lk.dir, name)
		f, err := os.Create(fname)
		if err == nil {
			symlink := filepath.Join(lk.dir, link)
			os.Remove(symlink)        // ignore err
			os.Symlink(name, symlink) // ignore err
			lk.add(tag, &fileBlock{fileInfo: fileInfo{name: name, timestamp: ""}, next: nil})
			return f, fname, nil
		}
		lastErr = err
	}
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}
