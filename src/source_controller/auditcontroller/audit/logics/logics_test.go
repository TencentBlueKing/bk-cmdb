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
 
package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/source_controller/common/commondata"
	storage "configcenter/src/storage"
	"errors"
	"testing"
)

func TestAddLogMutilNotData(t *testing.T) {

	mockdb := &mockMongo{
		data: nil,
		err:  nil,
	}
	DB = mockdb
	err := AddLogMulti(1, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, nil, "null", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}
}

func TestAddLogMutilErr(t *testing.T) {

	mockdb := &mockMongo{
		data:           nil,
		err:            errors.New("mock insert all"),
		errTrigger:     1,
		errTriggerStep: 0,
	}
	DB = mockdb

	contents := []auditoplog.AuditLogContext{
		auditoplog.AuditLogContext{ID: 1, Content: "sss"},
	}

	err := AddLogMulti(1, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, contents, "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}

}

func TestAddLogMutil(t *testing.T) {

	mockdb := &mockMongo{
		data: nil,
		err:  nil,
	}
	DB = mockdb

	contents := []auditoplog.AuditLogContext{
		auditoplog.AuditLogContext{ID: 1, Content: "sss"},
	}

	err := AddLogMulti(1, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, contents, "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}

}

func TestAddLogMultiWithExtKeyNotData(t *testing.T) {
	mockdb := &mockMongo{
		data: nil,
		err:  nil,
	}
	DB = mockdb

	err := AddLogMultiWithExtKey(1, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, nil, "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}
}

func TestAddLogMultiWithExtKeyErr(t *testing.T) {
	mockdb := &mockMongo{
		data:           nil,
		err:            errors.New("mock error, TestAddLogMultiWithExtKeyErr"),
		errTrigger:     1,
		errTriggerStep: 0,
	}
	DB = mockdb

	contents := []auditoplog.AuditLogExt{
		auditoplog.AuditLogExt{ID: 1, Content: "row1", ExtKey: "127.0.0.1"},
	}

	err := AddLogMultiWithExtKey(1, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, contents, "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}
}

func TestAddLogMultiWithExtKey(t *testing.T) {
	mockdb := &mockMongo{
		data: nil,
		err:  nil,
	}
	DB = mockdb

	contents := []auditoplog.AuditLogExt{
		auditoplog.AuditLogExt{ID: 1, Content: "row1", ExtKey: "127.0.0.1"},
		auditoplog.AuditLogExt{ID: 2, Content: "row2", ExtKey: "127.0.0.2"},
	}

	err := AddLogMultiWithExtKey(1, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, contents, "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}
}

func TestAddLogWithStr(t *testing.T) {
	mockdb := &mockMongo{
		data: nil,
		err:  nil,
	}
	DB = mockdb

	err := AddLogWithStr(1, 0, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, "test TestAddLogWithStr", "key", "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}
}

func TestAddLogWithStrErr(t *testing.T) {
	mockdb := &mockMongo{
		data:           nil,
		err:            errors.New("mock error, TestAddLogWithStrErr"),
		errTrigger:     1,
		errTriggerStep: 0,
	}
	DB = mockdb

	err := AddLogWithStr(1, 0, auditoplog.AuditOpTypeAdd, common.BKInnerObjIDHost, "test TestAddLogWithStr", "key", "mock desc", common.BKDefaultOwnerID, "user")
	if err != mockdb.err {
		t.Error(err)
	}
}

func TestSearchGetDataErr(t *testing.T) {
	mockdb := &mockMongo{
		data:           nil,
		err:            errors.New("mock error, TestSearchGetDataErr"),
		errTrigger:     1,
		errTriggerStep: 0,
	}
	DB = mockdb
	var dat commondata.ObjQueryInput
	_, _, err := Search(dat)

	if mockdb.err != err {
		t.Error(err)
	}

}

func TestSearchGetCntErr(t *testing.T) {
	mockdb := &mockMongo{
		data:           nil,
		err:            errors.New("mock error, TestSearchGetCntErr"),
		errTrigger:     2,
		errTriggerStep: 0,
	}
	DB = mockdb
	var dat commondata.ObjQueryInput
	_, _, err := Search(dat)

	if mockdb.err != err {
		t.Error(err)
	}

}

func TestSearch(t *testing.T) {
	mockdb := &mockMongo{
		data: nil,
		err:  nil,
	}
	DB = mockdb
	var dat commondata.ObjQueryInput
	_, _, err := Search(dat)

	if mockdb.err != err {
		t.Error(err)
	}

}

type mockMongo struct {
	data           []interface{}
	err            error
	errTrigger     int
	errTriggerStep int
}

func (m *mockMongo) Open() error {
	return nil
}

func (m *mockMongo) Insert(tableName string, data interface{}) (int, error) {
	m.errTriggerStep++
	if m.errTrigger == m.errTriggerStep {
		return 0, m.err
	}

	m.data = append(m.data, data)
	return 1, nil
}

func (m *mockMongo) InsertMuti(cName string, data ...interface{}) error {
	m.errTriggerStep++
	if m.errTrigger == m.errTriggerStep {
		return m.err
	}

	m.data = append(m.data, data...)
	return nil
}

func (m *mockMongo) GetIncID(cName string) (int64, error) {
	return 1, nil
}

func (m *mockMongo) UpdateByCondition(cName string, data, condiction interface{}) error {
	return nil
}
func (m *mockMongo) GetOneByCondition(cName string, fields []string, condiction interface{}, result interface{}) error {
	return nil
}
func (m *mockMongo) GetMutilByCondition(cName string, fields []string, condiction interface{}, result interface{}, sort string, start, limit int) error {
	m.errTriggerStep++
	if m.errTrigger == m.errTriggerStep {
		return m.err
	}
	return nil
}
func (m *mockMongo) GetCntByCondition(cName string, condiction interface{}) (int, error) {
	m.errTriggerStep++
	if m.errTrigger == m.errTriggerStep {
		return 0, m.err
	}
	return 1, nil
}
func (m *mockMongo) DelByCondition(cName string, condiction interface{}) error {
	return nil
}
func (m *mockMongo) HasTable(cName string) (bool, error) {
	return true, nil
}
func (m *mockMongo) ExecSql(cmd interface{}) error {
	return nil
}
func (m *mockMongo) Index(cName string, index *storage.Index) error {
	return nil
}
func (m *mockMongo) DropTable(cName string) error {
	return nil
}
func (m *mockMongo) HasFields(cName, field string) (bool, error) {
	return true, nil
}
func (m *mockMongo) AddColumn(cName string, column *storage.Column) error {
	return nil
}
func (m *mockMongo) ModifyColumn(cName, oldName, newColumn string) error {
	return nil
}
func (m *mockMongo) DropColumn(cName, field string) error {
	return nil
}
func (m *mockMongo) CreateTable(sql string) error {
	return nil
}
func (m *mockMongo) GetType() string {
	return "mock db"
}

func (m *mockMongo) Close() {
	return
}
func (m *mockMongo) GetSession() interface{} {
	return nil
}
