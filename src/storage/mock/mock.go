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
