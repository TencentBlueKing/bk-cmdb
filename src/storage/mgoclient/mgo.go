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

package mgoclient

import (
	"errors"
	"fmt"
	"strings"

	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage"
	// "log"
	// "os"
)

type MongoConfig struct {
	Address      string
	User         string
	Password     string
	Database     string
	Port         string
	MaxOpenConns string
	MaxIdleConns string
	Mechanism    string
}

func NewMongoConfig(src map[string]string) *MongoConfig {
	config := MongoConfig{}
	config.Address = src["mongodb.host"]
	config.User = src["mongodb.usr"]
	config.Password = src["mongodb.pwd"]
	config.Database = src["mongodb.database"]
	config.Port = src["mongodb.port"]
	config.MaxOpenConns = src["mongodb.maxOpenConns"]
	config.MaxIdleConns = src["mongodb.maxIDleConns"]
	return &config
}

type MgoCli struct {
	host      string
	port      string
	usr       string
	pwd       string
	dbName    string
	mechanism string
	session   *mgo.Session
}

func NewFromConfig(cfg MongoConfig) (*MgoCli, error) {
	return NewMgoCli(cfg.Address, cfg.Port, cfg.User, cfg.Password, cfg.Mechanism, cfg.Database)
}

func NewMgoCli(host, port, usr, pwd, mechanism, database string) (*MgoCli, error) {
	mgocli := new(MgoCli)
	mgocli.host = host
	mgocli.port = port
	mgocli.usr = usr
	mgocli.pwd = pwd
	mgocli.dbName = database
	mgocli.mechanism = mechanism
	return mgocli, nil
}

// Open open the connection
func (m *MgoCli) Open() error {

	dialInfo := &mgo.DialInfo{
		Addrs:     strings.Split(m.host, ","),
		Direct:    false,
		Timeout:   time.Second * 5,
		Database:  m.dbName,
		Source:    "",
		Username:  m.usr,
		Password:  m.pwd,
		PoolLimit: 4096,
		Mechanism: m.mechanism,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	m.session = session
	if err != nil {
		return err
	}
	return nil
}

// Ping will ping the db server
func (m *MgoCli) Ping() error {
	return m.session.Ping()
}

// GetSession returns mongo session
func (m *MgoCli) GetSession() interface{} {
	return m.session
}

// Close close mongo session
func (m *MgoCli) Close() {
	if m.session != nil {
		m.session.Close()
	}
}

// Insert insert one document
func (m *MgoCli) Insert(cName string, data interface{}) (int, error) {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C(cName)
	EscapeHtml(data)
	err := c.Insert(data)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// InsertMuti insert muti documents
func (m *MgoCli) InsertMuti(cName string, data ...interface{}) error {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C(cName)
	EscapeHtml(data...)
	err := c.Insert(data...)
	if err != nil {
		return err
	}
	return nil
}

// UpdateByCondition update documents by condiction
func (m *MgoCli) UpdateByCondition(cName string, data, condition interface{}) error {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C(cName)
	EscapeHtml(data)
	datac := bson.M{"$set": data}
	_, err := c.UpdateAll(condition, datac)
	if err != nil {
		return err
	}
	return nil
}

// GetOneByCondition get one document by condiction
func (m *MgoCli) GetOneByCondition(cName string, fields []string, condiction interface{}, result interface{}) error {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C(cName)
	fieldmap := make(map[string]interface{})
	if 0 != len(fields) {
		for _, key := range fields {
			fieldmap[key] = 1
		}
	}

	fieldmap["_id"] = 0
	query := c.Find(condiction)
	if 0 < len(fieldmap) {
		query.Select(fieldmap)
	}
	return query.One(result)
}

// GetMutilByCondition get multiple document by condiction
func (m *MgoCli) GetMutilByCondition(cName string, fields []string, condiction interface{}, result interface{}, sort string, start, limit int) error {
	m.session.Refresh()
	if len(fields) == 1 && fields[0] == "" {
		fields = nil
	}
	c := m.session.DB(m.dbName).C(cName)
	fieldmap := make(map[string]interface{})
	if 0 != len(fields) {
		for _, key := range fields {
			fieldmap[key] = 1
		}
	}

	fieldmap["_id"] = 0
	query := c.Find(condiction)
	if 0 < len(fieldmap) {
		query = query.Select(fieldmap)
	}
	if "" != sort {
		arrSort := strings.Split(sort, common.BKDBSortFieldSep)
		query = query.Sort(arrSort...)
	}

	if 0 < start {
		query = query.Skip(start)
	}
	if 0 < limit {
		query = query.Limit(limit)
	}
	err := query.All(result)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return err
	}
	return nil
}

// GetCntByCondition returns count number filter by condiction
func (m *MgoCli) GetCntByCondition(cName string, condition interface{}) (cnt int, err error) {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C(cName)
	count := 0
	count, err = c.Find(condition).Count()
	if err != nil {
		return count, err
	}
	return count, nil
}

// GetIncID returns next sequence ID for cName collection
//
// db.cc_idgenerator.findAndModify(
// {
// 	query:{_id: "sub" },
// 	update: {$inc:{SequenceID:1}},
// 	upsert:true,
// 	new:true
//  }).sequence_value
func (m *MgoCli) GetIncID(cName string) (incID int64, err error) {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C("cc_idgenerator")
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"SequenceID": int64(1)}},
		ReturnNew: true,
		Upsert:    true,
	}
	doc := map[string]interface{}{}
	_, err = c.Find(bson.M{"_id": cName}).Apply(change, &doc)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(fmt.Sprint(doc["SequenceID"]), 10, 64)
}

//按条件删除主句
func (m *MgoCli) DelByCondition(cName string, condiction interface{}) error {
	m.session.Refresh()
	c := m.session.DB(m.dbName).C(cName)

	// s = s.(map[string]interface{})
	// selector := bson.M(s)
	_, err := c.RemoveAll(condiction)
	if err != nil {
		return err
	}
	return nil
}

//判断表是否存在
func (m *MgoCli) HasTable(tableName string) (bool, error) {
	m.session.Refresh()
	tableNames, err := m.session.DB(m.dbName).CollectionNames()
	if nil != err {
		blog.Error("mongo query failure, %v", err)
		return false, err
	}
	for _, tn := range tableNames {
		if tn == tableName {
			return true, nil
		}
	}
	return false, nil
}

func (m *MgoCli) CreateTable(tableName string) error {
	m.session.Refresh()
	cInfo := new(mgo.CollectionInfo)
	c := m.session.DB(m.dbName).C(tableName)
	err := c.Create(cInfo)
	if nil != err {
		return err
	}
	return nil
}

//执行原始的sql语句，并不会返回数据， 只会是否执行出错
func (m *MgoCli) ExecSql(cmd interface{}) error {
	return errors.New("not support method")
}

func (m *MgoCli) Index(tableName string, index *storage.Index) error {
	m.session.Refresh()
	unique := false
	backgroud := false
	switch index.Type {
	case storage.INDEX_TYPE_BACKGROUP_UNIQUE:
		unique = true
		backgroud = true
	case storage.INDEX_TYPE_UNIQUE:
		unique = true
	case storage.INDEX_TYPE_BACKGROUP:
		backgroud = true
	}
	return m.session.DB(m.dbName).C(tableName).EnsureIndex(mgo.Index{
		Name:       index.Name,
		Key:        index.Columns,
		Unique:     unique,
		Background: backgroud,
	})
}

func (m *MgoCli) DropTable(tableName string) error {
	m.session.Refresh()
	return m.session.DB(m.dbName).C(tableName).DropCollection()
}

func (m *MgoCli) HasFields(tableName, field string) (bool, error) {
	m.session.Refresh()
	selector := bson.M{field: bson.M{"$exists": true}}
	count, err := m.session.DB(m.dbName).C(tableName).Find(selector).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//新加字段， 表名，字段名,字段类型, 附加描述（是否为空， 默认值）
func (m *MgoCli) AddColumn(tableName string, column *storage.Column) error {
	m.session.Refresh()
	selector := bson.M{column.Name: bson.M{"$exists": false}}
	datac := bson.M{"$set": bson.M{column.Name: column.Ext}}
	_, err := m.session.DB(m.dbName).C(tableName).UpdateAll(selector, datac)
	return err
}

func (m *MgoCli) ModifyColumn(tableName, oldName, newColumn string) error {
	m.session.Refresh()
	datac := bson.M{"$rename": bson.M{oldName: newColumn}}
	_, err := m.session.DB(m.dbName).C(tableName).UpdateAll(nil, datac)
	return err
}

func (m *MgoCli) DropColumn(tableName, field string) error {
	m.session.Refresh()
	datac := bson.M{"$unset": bson.M{field: "1"}}
	_, err := m.session.DB(m.dbName).C(tableName).UpdateAll(nil, datac)
	return err
}

//GetType 获取操作db的类
func (m *MgoCli) GetType() string {
	return storage.DI_MONGO
}

// IsDuplicateErr returns whether err is duplicate error
func (m *MgoCli) IsDuplicateErr(err error) bool {
	return mgo.IsDup(err)
}

// IsNotFoundErr returns whether err is not found error
func (m *MgoCli) IsNotFoundErr(err error) bool {
	return mgo.ErrNotFound == err
}
func (m *MgoCli) GetDBName() string {
	return m.dbName
}
