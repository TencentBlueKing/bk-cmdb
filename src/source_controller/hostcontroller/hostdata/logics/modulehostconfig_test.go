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
    "configcenter/src/source_controller/common/eventdata"
    "configcenter/src/source_controller/common/instdata"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/storage"
    "testing"
    "errors"
    "flag"
)

type MockDI struct {
    ErrGetIncID            error
    ErrInsert              error
    ErrInsertMuti          error
    ErrUpdateByCondition   error
    ErrGetOneByCondition   error
    ErrGetMutilByCondition error
    ErrGetCntByCondition   error
    ErrDelByCondition      error
    ErrHasTable            error
    ErrExecSql             error
    ErrIndex               error
    ErrDropTable           error
    ErrHasFields           error
    ErrAddColumn           error
    ErrModifyColumn        error
    ErrDropColumn          error
    ErrCreateTable         error
    ErrOpen                error
    ErrGetSession          error

    VarGetIncID          int64
    VarInsert            int
    VarGetCntByCondition int
    VarHasTable          bool
    VarHasFields         bool
    VarGetType           string
}

func (m *MockDI) GetIncID(cName string) (int64, error) {return m.VarGetIncID, m.ErrGetIncID}
func (m *MockDI) Insert(cName string, data interface{}) (int, error) {return m.VarInsert, m.ErrInsert}
func (m *MockDI) InsertMuti(cName string, data ...interface{}) error {return m.ErrInsertMuti}
func (m *MockDI) UpdateByCondition(cName string, data, condition interface{}) error {return m.ErrUpdateByCondition}
func (m *MockDI) GetOneByCondition(cName string, fields []string, condition interface{}, result interface{}) error {return m.ErrGetOneByCondition}
func (m *MockDI) GetMutilByCondition(cName string, fields []string, condition interface{}, result interface{}, sort string, start, limit int) error {return m.ErrGetMutilByCondition}
func (m *MockDI) GetCntByCondition(cName string, condition interface{}) (int, error) {return m.VarGetCntByCondition, m.ErrGetCntByCondition}
func (m *MockDI) DelByCondition(cName string, condition interface{}) error {return m.ErrDelByCondition}
func (m *MockDI) HasTable(cName string) (bool, error) {return m.VarHasTable, m.ErrHasTable}
func (m *MockDI) ExecSql(cmd interface{}) error {return m.ErrExecSql}
func (m *MockDI) Index(cName string, index *storage.Index) error {return m.ErrIndex}
func (m *MockDI) DropTable(cName string) error {return m.ErrDropTable}
func (m *MockDI) HasFields(cName, field string) (bool, error) {return m.VarHasFields, m.ErrHasFields}
func (m *MockDI) AddColumn(cName string, column *storage.Column) error {return m.ErrAddColumn}
func (m *MockDI) ModifyColumn(cName, oldName, newColumn string) error {return m.ErrModifyColumn}
func (m *MockDI) DropColumn(cName, field string) error {return m.ErrDropColumn}
func (m *MockDI) CreateTable(sql string) error {return m.ErrCreateTable}
func (m *MockDI) GetType() string {return m.VarGetType}
func (m *MockDI) Open() error {return m.ErrOpen}
func (m *MockDI) Close() {}
func (m *MockDI) GetSession() interface{} {return m.ErrGetSession}

func TestDelSingleHostModuleRelation(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    instdata.DataH = &MockDI{ErrGetOneByCondition: errFake}
    _, err := DelSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if err != errFake {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestDelSingleHostModuleRelation2(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetCntByCondition: errFake}
    instdata.DataH = &MockDI{}
    _, err := DelSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if err != errFake {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestDelSingleHostModuleRelation3(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{VarGetCntByCondition: 0}
    instdata.DataH = &MockDI{}
    r, err := DelSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if !r {
        t.Errorf("result not as expected, should be true")
    }
    if err != nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestDelSingleHostModuleRelation4(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetMutilByCondition: errFake, VarGetCntByCondition: 1}
    instdata.DataH = &MockDI{}
    _, err := DelSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if err != errFake {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestDelSingleHostModuleRelation5(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrDelByCondition: errFake, VarGetCntByCondition: 1}
    instdata.DataH = &MockDI{}
    _, err := DelSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if err != errFake {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestDelSingleHostModuleRelation6(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{VarGetCntByCondition: 1}
    instdata.DataH = &MockDI{}
    r, err := DelSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if !r {
        t.Errorf("result not as expected, should be true")
    }
    if err != nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestAddSingleHostModuleRelation(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    instdata.DataH = &MockDI{ErrGetOneByCondition: errFake}
    _, err := AddSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if err != errFake {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestAddSingleHostModuleRelation2(t *testing.T) {
    ec := &eventdata.EventContext{}
    cc := &api.APIResource{}

    instdata.DataH = &MockDI{}
    r, err := AddSingleHostModuleRelation(ec, cc, 1, 2, 3)
    if r {
        t.Errorf("result not as expected, should be false")
    }
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetDefaultModuleIDs(t *testing.T) {
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetMutilByCondition: errFake}
    _, err := GetDefaultModuleIDs(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetDefaultModuleIDs2(t *testing.T) {
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{}
    _, err := GetDefaultModuleIDs(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetModuleIDsByHostID(t *testing.T) {
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetMutilByCondition: errFake}
    _, err := GetModuleIDsByHostID(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetModuleIDsByHostID2(t *testing.T) {
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{}
    _, err := GetModuleIDsByHostID(cc, 1)
    if err != nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetResourcePoolApp(t *testing.T) {
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetOneByCondition: errFake}
    _, err := GetResourcePoolApp(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetResourcePoolApp2(t *testing.T) {
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{}
    _, err := GetResourcePoolApp(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestCheckHostInIDle(t *testing.T) {
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetMutilByCondition: errFake}
    _, err := CheckHostInIDle(cc, 1, 2, nil)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestCheckHostInIDle2(t *testing.T) {
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{}
    _, err := CheckHostInIDle(cc, 1, 2, nil)
    if err != nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetIDleModuleID(t *testing.T) {
    cc := &api.APIResource{}

    errFake := errors.New("fake error")
    cc.InstCli = &MockDI{ErrGetOneByCondition: errFake}
    _, err := GetIDleModuleID(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func TestGetIDleModuleID2(t *testing.T) {
    cc := &api.APIResource{}

    cc.InstCli = &MockDI{}
    _, err := GetIDleModuleID(cc, 1)
    if err == nil {
        t.Errorf("error not as expected: %v", err)
    }
}

func init() {
    flag.Set("logtostderr", "false")
}
