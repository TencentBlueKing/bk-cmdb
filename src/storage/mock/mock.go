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

package mock

import (
	"configcenter/src/storage"
)

type MockDB struct{}

func (m *MockDB) GetIncID(cName string) (int64, error) {
	return 0, nil
}
func (m *MockDB) Insert(cName string, data interface{}) (int, error) {
	return 0, nil
}
func (m *MockDB) InsertMuti(cName string, data ...interface{}) error                { return nil }
func (m *MockDB) IsDuplicateErr(err error) bool                                     { return false }
func (m *MockDB) IsNotFoundErr(err error) bool                                      { return false }
func (m *MockDB) UpdateByCondition(cName string, data, condition interface{}) error { return nil }
func (m *MockDB) GetOneByCondition(cName string, fields []string, condition interface{}, result interface{}) error {
	return nil
}
func (m *MockDB) GetMutilByCondition(cName string, fields []string, condition interface{}, result interface{}, sort string, start, limit int) error {
	return nil
}
func (m *MockDB) GetCntByCondition(cName string, condition interface{}) (int, error) { return 0, nil }
func (m *MockDB) DelByCondition(cName string, condition interface{}) error           { return nil }
func (m *MockDB) HasTable(cName string) (bool, error)                                { return false, nil }
func (m *MockDB) ExecSql(cmd interface{}) error                                      { return nil }
func (m *MockDB) Index(cName string, index *storage.Index) error                     { return nil }
func (m *MockDB) DropTable(cName string) error                                       { return nil }
func (m *MockDB) HasFields(cName, field string) (bool, error)                        { return false, nil }
func (m *MockDB) AddColumn(cName string, column *storage.Column) error               { return nil }
func (m *MockDB) ModifyColumn(cName, oldName, newColumn string) error                { return nil }
func (m *MockDB) DropColumn(cName, field string) error                               { return nil }
func (m *MockDB) CreateTable(sql string) error                                       { return nil }
func (m *MockDB) GetType() string                                                    { return "" }
func (m *MockDB) Open() error                                                        { return nil }
func (m *MockDB) Ping() error                                                        { return nil }
func (m *MockDB) Close()                                                             {}
func (m *MockDB) GetSession() interface{}                                            { return nil }
func (m *MockDB) GetDBName() string                                                  { return "mock_db" }
