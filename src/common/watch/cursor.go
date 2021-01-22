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

package watch

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// this cursor means there is no event occurs.
// we send this cursor to our the watcher and if we
// received a NoEventCursor, then we need to fetch event
// from the head cursor
var NoEventCursor string

func init() {
	no := Cursor{
		Type:        NoEvent,
		ClusterTime: types.TimeStamp{Sec: 1, Nano: 1},
		Oid:         "5ea6d3f394c1f5d986e9bd86",
	}
	cursor, err := no.Encode()
	if err != nil {
		panic("initial NoEventCursor failed")
	}
	// cursor should be:
	// MQ0xDTVlYTZkM2YzOTRjMWY1ZDk4NmU5YmQ4Ng0xDTE=
	NoEventCursor = cursor
}

type CursorType string

const (
	NoEvent                 CursorType = "no_event"
	UnknownType             CursorType = "unknown"
	Host                    CursorType = "host"
	ModuleHostRelation      CursorType = "host_relation"
	Biz                     CursorType = "biz"
	Set                     CursorType = "set"
	Module                  CursorType = "module"
	SetTemplate             CursorType = "set_template"
	ObjectBase              CursorType = "object_instance"
	Process                 CursorType = "process"
	ProcessInstanceRelation CursorType = "process_instance_relation"
)

func (ct CursorType) ToInt() int {
	switch ct {
	case NoEvent:
		return 1
	case Host:
		return 2
	case ModuleHostRelation:
		return 3
	case Biz:
		return 4
	case Set:
		return 5
	case Module:
		return 6
	case SetTemplate:
		return 7
	case ObjectBase:
		return 8
	case Process:
		return 9
	case ProcessInstanceRelation:
		return 10
	default:
		return -1
	}
}

func (ct *CursorType) ParseInt(typ int) {
	switch typ {
	case 1:
		*ct = NoEvent
	case 2:
		*ct = Host
	case 3:
		*ct = ModuleHostRelation
	case 4:
		*ct = Biz
	case 5:
		*ct = Set
	case 6:
		*ct = Module
	case 7:
		*ct = SetTemplate
	case 8:
		*ct = ObjectBase
	case 9:
		*ct = Process
	case 10:
		*ct = ProcessInstanceRelation
	default:
		*ct = UnknownType
	}
}

// ListCursorTypes returns all support CursorTypes.
func ListCursorTypes() []CursorType {
	return []CursorType{Host, ModuleHostRelation, Biz, Set, Module, SetTemplate, ObjectBase, Process, ProcessInstanceRelation}
}

// ListEventCallbackCursorTypes returns all support CursorTypes for event callback.
func ListEventCallbackCursorTypes() []CursorType {
	return []CursorType{Host, ModuleHostRelation, Biz, Set, Module, ObjectBase, Process, ProcessInstanceRelation}
}

// Cursor is a self-defined token which is corresponding to the mongodb's resume token.
// cursor has a unique and 1:1 relationship with mongodb's resume token.
type Cursor struct {
	Type        CursorType
	ClusterTime types.TimeStamp
	// a random hex string to avoid the caller to generated a self-defined cursor.
	Oid string
}

const cursorVersion = "1"

func (c Cursor) Encode() (string, error) {
	if c.Type == "" {
		return "", errors.New("unsupported type")
	}

	if c.ClusterTime.Sec == 0 {
		return "", errors.New("invalid cluster time sec")
	}

	if c.Oid == "" {
		return "", errors.New("invalid oid")
	}

	sec := strconv.FormatUint(uint64(c.ClusterTime.Sec), 10)
	nano := strconv.FormatUint(uint64(c.ClusterTime.Nano), 10)
	pool := bytes.Buffer{}
	// version field.
	pool.WriteString(cursorVersion)
	pool.WriteByte('\r')

	// type filed.
	if c.Type.ToInt() < 0 {
		return "", errors.New("unsupported cursor type")
	}

	pool.WriteString(strconv.Itoa(c.Type.ToInt()))
	pool.WriteByte('\r')

	// oid field.
	pool.WriteString(c.Oid)
	pool.WriteByte('\r')

	// cluster time sec field.
	pool.WriteString(sec)
	pool.WriteByte('\r')

	// cluster time nano field
	pool.WriteString(nano)

	return base64.StdEncoding.EncodeToString(pool.Bytes()), nil
}

func (c *Cursor) Decode(cur string) error {
	byt, err := base64.StdEncoding.DecodeString(cur)
	if err != nil {
		return fmt.Errorf("decode cursor, but base64 decode failed, err: %v", err)
	}

	elements := make([]string, 0)
	pool := bytes.NewBuffer(byt)

	ele := make([]byte, 0)
	for {
		b, err := pool.ReadByte()
		if err != nil {
			if err != io.EOF {
				return err
			}
			// to the end
			elements = append(elements, string(ele))
			break
		}
		if b == '\r' {
			elements = append(elements, string(ele))
			ele = ele[:0]
		} else {
			ele = append(ele, b)
		}
	}

	if len(elements) != 5 {
		return errors.New("invalid cursor string")
	}

	if elements[0] != cursorVersion {
		return fmt.Errorf("decode cursor, but got invalid cursor version: %s", elements[0])
	}

	typ, err := strconv.Atoi(elements[1])
	if err != nil {
		return fmt.Errorf("got invalid type: %s", elements[1])
	}
	cursorType := CursorType("")
	cursorType.ParseInt(typ)
	c.Type = cursorType

	_, err = primitive.ObjectIDFromHex(elements[2])
	if err != nil {
		return fmt.Errorf("got invalid oid: %s, err: %v", elements[2], err)
	}
	c.Oid = elements[2]

	sec, err := strconv.ParseUint(elements[3], 10, 64)
	if err != nil {
		return fmt.Errorf("got invalid sec field %s, err: %v", elements[3], err)
	}
	c.ClusterTime.Sec = uint32(sec)

	nano, err := strconv.ParseUint(elements[4], 10, 64)
	if err != nil {
		return fmt.Errorf("got invalid nano field %s, err: %v", elements[4], err)
	}
	c.ClusterTime.Nano = uint32(nano)
	return nil
}

func GetEventCursor(coll string, e *types.Event) (string, error) {
	curType := UnknownType
	switch coll {
	case common.BKTableNameBaseHost:
		curType = Host
	case common.BKTableNameModuleHostConfig:
		curType = ModuleHostRelation
	case common.BKTableNameBaseApp:
		curType = Biz
	case common.BKTableNameBaseSet:
		curType = Set
	case common.BKTableNameBaseModule:
		curType = Module
	case common.BKTableNameSetTemplate:
		curType = SetTemplate
	case common.BKTableNameBaseInst:
		curType = ObjectBase
	case common.BKTableNameBaseProcess:
		curType = Process
	case common.BKTableNameProcessInstanceRelation:
		curType = ProcessInstanceRelation
	default:
		blog.Errorf("unsupported cursor type collection: %s, oid: %s", e.Oid)
		return "", fmt.Errorf("unsupported cursor type collection: %s", coll)
	}

	hCursor := &Cursor{
		Type:        curType,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
	}

	hCursorEncode, err := hCursor.Encode()
	if err != nil {
		blog.Errorf("encode head node cursor failed, err: %v, oid: %s", err, e.Oid)
		return "", err
	}

	return hCursorEncode, nil
}
