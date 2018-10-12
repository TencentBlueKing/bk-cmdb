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

package storage

// DI define storage interface
type DI interface {
	GetIncID(cName string) (int64, error)
	Insert(cName string, data interface{}) (int, error)
	InsertMuti(cName string, data ...interface{}) error
	IsDuplicateErr(err error) bool
	IsNotFoundErr(err error) bool
	UpdateByCondition(cName string, data, condition interface{}) error
	GetOneByCondition(cName string, fields []string, condition interface{}, result interface{}) error
	GetMutilByCondition(cName string, fields []string, condition interface{}, result interface{}, sort string, start, limit int) error
	GetCntByCondition(cName string, condition interface{}) (int, error)
	DelByCondition(cName string, condition interface{}) error
	HasTable(cName string) (bool, error)
	ExecSql(cmd interface{}) error
	Index(cName string, index *Index) error
	DropTable(cName string) error
	HasFields(cName, field string) (bool, error)
	AddColumn(cName string, column *Column) error
	ModifyColumn(cName, oldName, newColumn string) error
	DropColumn(cName, field string) error
	CreateTable(sql string) error
	GetType() string
	Open() error
	Ping() error
	Close()
	GetSession() interface{}
	GetDBName() string
}

const (
	INDEX_TYPE_UNIQUE           = 1 //唯一索引
	INDEX_TYPE_PRIMAEY          = 2 //主键
	INDEX_TYPE_GENERAL          = 3 // 普通索引
	INDEX_TYPE_BACKGROUP_UNIQUE = 4 //mogon特殊，唯一，且后台生产索引
	INDEX_TYPE_BACKGROUP        = 5 //mogon特殊，后台生产索引
)

type Index struct {
	Name    string
	Columns []string
	Type    int
}

type Column struct {
	Name string
	Ext  interface{}
}

type M map[string]interface{}

func GetMongoIndex(name string, keys []string, unique, background bool) *Index {
	indexType := INDEX_TYPE_GENERAL
	if true == unique && false == background {
		indexType = INDEX_TYPE_BACKGROUP_UNIQUE
	} else if true == unique {
		indexType = INDEX_TYPE_UNIQUE
	} else if true == background {
		indexType = INDEX_TYPE_BACKGROUP
	}

	return &Index{
		Name:    name,
		Columns: keys,
		Type:    indexType,
	}
}

func GetMongoColumn(name string, val interface{}) *Column {
	return &Column{
		Name: name,
		Ext:  val,
	}
}

const (
	DI_MYSQL string = "mysql"
	DI_MONGO string = "mongodb"
	DI_REDIS string = "redis"
)
