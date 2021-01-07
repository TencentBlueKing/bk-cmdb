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

package configcenter

import (
	"bytes"
	err "errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/viper"
)

// configuration parser representing five files,the configuration can be taken out according to the parser.
var redisParser *viperParser
var mongodbParser *viperParser
var commonParser *viperParser
var extraParser *viperParser
var migrateParser *viperParser

var confLock sync.RWMutex

func checkDir(path string) error {
	info, err := os.Stat(path)
	if os.ErrNotExist == err {
		return fmt.Errorf("directory %s not exists", path)
	}
	if err != nil {
		return fmt.Errorf("stat directory %s faile, %s", path, err.Error())
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not directory", path)
	}

	return nil
}

func loadErrorAndLanguage(errorres string, languageres string, handler *CCHandler) error {
	if err := checkDir(errorres); err != nil {
		return err
	}
	errcode, err := errors.LoadErrorResourceFromDir(errorres)
	if err != nil {
		return fmt.Errorf("load error resource error: %s", err)
	}
	handler.OnErrorUpdate(nil, errcode)

	if err := checkDir(languageres); err != nil {
		return err
	}
	languagepack, err := language.LoadLanguageResourceFromDir(languageres)
	if err != nil {
		return fmt.Errorf("load language resource error: %s", err)
	}
	handler.OnLanguageUpdate(nil, languagepack)
	return nil
}

func LoadConfigFromLocalFile(confPath string, handler *CCHandler) error {

	// if it is admin_server, skip the loading of other files,load only error and language.
	if common.GetIdentification() == types.CC_MODULE_MIGRATE {
		errorres, _ := String("errors.res")
		languageres, _ := String("language.res")
		return loadErrorAndLanguage(errorres, languageres, handler)
	}

	// load local error and language
	errorres := confPath + "/errors"
	languageres := confPath + "/language"
	if err := loadErrorAndLanguage(errorres, languageres, handler); err != nil {
		return err
	}

	// load local common
	commonPath := confPath + "/" + types.CCConfigureCommon
	if err := SetCommonFromFile(commonPath); err != nil {
		blog.Errorf("load config from file[%s], but can not found common config", commonPath)
		return err
	}
	if handler.OnProcessUpdate != nil {
		handler.OnProcessUpdate(ProcessConfig{}, ProcessConfig{})
	}

	// load local extra
	extraPath := confPath + "/" + types.CCConfigureExtra
	if err := SetExtraFromFile(extraPath); err != nil {
		blog.Errorf("load config from file[%s], but can not found extra config", extraPath)
		return err
	}
	if handler.OnExtraUpdate != nil {
		handler.OnExtraUpdate(ProcessConfig{}, ProcessConfig{})
	}

	// load local redis
	redisPath := confPath + "/" + types.CCConfigureRedis
	if err := SetRedisFromFile(redisPath); err != nil {
		return err
	}

	// load local mongodb
	mongodbPath := confPath + "/" + types.CCConfigureMongo
	if err := SetMongodbFromFile(mongodbPath); err != nil {
		return err
	}

	return nil
}

func SetRedisFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	if redisParser != nil {
		err := redisParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from redis")
			return err
		}
		return nil
	}
	redisParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from redis")
		return err
	}
	return nil
}

func SetRedisFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	redisParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from redis")
		return err
	}
	return nil
}

func SetMongodbFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	if mongodbParser != nil {
		err = mongodbParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from mongodb")
			return err
		}
		return nil
	}
	mongodbParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from mongodb")
		return err
	}
	return nil
}

func SetMongodbFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	mongodbParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from mongodb")
		return err
	}
	return nil
}

func SetCommonFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// if it is not nil, do not create a new parser, but add the new configuration information to viper
	if commonParser != nil {
		err = commonParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from common")
			return err
		}
		return nil
	}
	commonParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from common")
		return err
	}
	return nil
}

func SetCommonFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	commonParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from common")
		return err
	}
	return nil
}

func SetExtraFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	//if it is not nil, do not create a new parser, but add the new configuration information to viper
	if extraParser != nil {
		err = extraParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from extra")
			return err
		}
		return nil
	}
	extraParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from extra")
		return err
	}
	return nil
}

func SetExtraFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	extraParser, err = newViperParserFromFile(target)
	if err != nil {
		blog.Errorf("fail to read configure from extra")
		return err
	}
	return nil
}

func SetMigrateFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	if migrateParser != nil {
		err := migrateParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			blog.Errorf("fail to read configure from migrate")
			return err
		}
		return nil
	}
	migrateParser, err = newViperParser(data)
	if err != nil {
		blog.Errorf("fail to read configure from migrate")
		return err
	}
	return nil
}

func SetMigrateFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// /data/migrate.yaml -> /data/migrate
	split := strings.Split(target, ".")
	filePath := split[0]
	migrateParser, err = newViperParserFromFile(filePath)
	if err != nil {
		blog.Errorf("fail to read configure from migrate")
		return err
	}
	return nil
}

// Redis return redis configuration information according to the prefix.
func Redis(prefix string) (redis.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	var parser *viperParser
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		parser = getRedisParser()
		if parser != nil {
			break
		}
		blog.Warn("the configuration of redis is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		blog.Errorf("can't find redis configuration")
		return redis.Config{}, err.New("can't find redis configuration")
	}

	return redis.Config{
		Address:          parser.getString(prefix + ".host"),
		Password:         parser.getString(prefix + ".pwd"),
		Database:         parser.getString(prefix + ".database"),
		MasterName:       parser.getString(prefix + ".masterName"),
		SentinelPassword: parser.getString(prefix + ".sentinelPwd"),
		Enable:           parser.getString(prefix + ".enable"),
		MaxOpenConns:     parser.getInt(prefix + ".maxOpenConns"),
	}, nil
}

// Mongo return mongo configuration information according to the prefix.
func Mongo(prefix string) (mongo.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	var parser *viperParser
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		parser = getMongodbParser()
		if parser != nil {
			break
		}
		blog.Warn("the configuration of mongo is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		blog.Errorf("can't find mongo configuration")
		return mongo.Config{}, err.New("can't find mongo configuration")
	}

	c := mongo.Config{
		Address:   parser.getString(prefix + ".host"),
		Port:      parser.getString(prefix + ".port"),
		User:      parser.getString(prefix + ".usr"),
		Password:  parser.getString(prefix + ".pwd"),
		Database:  parser.getString(prefix + ".database"),
		Mechanism: parser.getString(prefix + ".mechanism"),
		RsName:    parser.getString(prefix + ".rsName"),
	}

	if c.RsName == "" {
		blog.Errorf("rsName not set")
	}
	if c.Mechanism == "" {
		c.Mechanism = "SCRAM-SHA-1"
	}
	if !parser.isSet(prefix+".maxOpenConns") || parser.getUint64(prefix+".maxOpenConns") > mongo.MaximumMaxOpenConns {
		c.MaxOpenConns = mongo.DefaultMaxOpenConns
	} else {
		c.MaxOpenConns = parser.getUint64(prefix + ".maxOpenConns")
	}

	if !parser.isSet(prefix+".maxIdleConns") || parser.getUint64(prefix+".maxIdleConns") < mongo.MinimumMaxIdleOpenConns {
		c.MaxIdleConns = mongo.MinimumMaxIdleOpenConns
	} else {
		c.MaxIdleConns = parser.getUint64(prefix + ".maxIdleConns")
	}

	if !parser.isSet(prefix + ".socketTimeoutSeconds") {
		blog.Errorf("can not find mongo.socketTimeoutSeconds config, use default value: %d", mongo.DefaultSocketTimeout)
		c.SocketTimeout = mongo.DefaultSocketTimeout
		return c, nil
	}

	c.SocketTimeout = parser.getInt(prefix + ".socketTimeoutSeconds")
	if c.SocketTimeout > mongo.MaximumSocketTimeout {
		blog.Errorf("mongo.socketTimeoutSeconds config %d exceeds maximum value, use maximum value %d", c.SocketTimeout, mongo.MaximumSocketTimeout)
		c.SocketTimeout = mongo.MaximumSocketTimeout
	}

	if c.SocketTimeout < mongo.MinimumSocketTimeout {
		blog.Errorf("mongo.socketTimeoutSeconds config %d less than minimum value, use minimum value %d", c.SocketTimeout, mongo.MinimumSocketTimeout)
		c.SocketTimeout = mongo.MinimumSocketTimeout
	}

	return c, nil
}

// String return the string value of the configuration information according to the key.
func String(key string) (string, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	if migrateParser != nil && migrateParser.isSet(key) {
		return migrateParser.getString(key), nil
	}
	if commonParser != nil && commonParser.isSet(key) {
		return commonParser.getString(key), nil
	}
	if extraParser != nil && extraParser.isSet(key) {
		return extraParser.getString(key), nil
	}
	return "", err.New("config not found")
}

// Int return the int value of the configuration information according to the key.
func Int(key string) (int, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	if migrateParser != nil && migrateParser.isSet(key) {
		return migrateParser.getInt(key), nil
	}
	if commonParser != nil && commonParser.isSet(key) {
		return commonParser.getInt(key), nil
	}
	if extraParser != nil && extraParser.isSet(key) {
		return extraParser.getInt(key), nil
	}
	return 0, err.New("config not found")
}

// Bool return the bool value of the configuration information according to the key.
func Bool(key string) (bool, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	if migrateParser != nil && migrateParser.isSet(key) {
		return migrateParser.getBool(key), nil
	}
	if commonParser != nil && commonParser.isSet(key) {
		return commonParser.getBool(key), nil
	}
	if extraParser != nil && extraParser.isSet(key) {
		return extraParser.getBool(key), nil
	}
	return false, err.New("config not found")
}

func IsExist(key string) bool {
	confLock.RLock()
	defer confLock.RUnlock()
	if migrateParser != nil {
		return migrateParser.isSet(key)
	}
	if (commonParser == nil || !commonParser.isSet(key)) && (extraParser == nil || !extraParser.isSet(key)) {
		return false
	}
	return true
}

func getRedisParser() *viperParser {
	if migrateParser != nil {
		return migrateParser
	}
	return redisParser
}

func getMongodbParser() *viperParser {
	if migrateParser != nil {
		return migrateParser
	}
	return mongodbParser
}

type Parser interface {
	GetString(string) string
	GetInt(string) int
	GetUint64(string) uint64
	GetBool(string) bool
	IsSet(path string) bool
}

type viperParser struct {
	parser *viper.Viper
}

func newViperParser(data []byte) (*viperParser, error) {

	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return &viperParser{parser: v}, nil
}

func newViperParserFromFile(target string) (*viperParser, error) {
	v := viper.New()
	v.SetConfigName(path.Base(target))
	v.AddConfigPath(path.Dir(target))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v.WatchConfig()
	return &viperParser{parser: v}, nil
}

func (vp *viperParser) getString(path string) string {
	return vp.parser.GetString(path)
}

func (vp *viperParser) getInt(path string) int {
	return vp.parser.GetInt(path)
}

func (vp *viperParser) getUint64(path string) uint64 {
	return vp.parser.GetUint64(path)
}

func (vp *viperParser) getBool(path string) bool {
	return vp.parser.GetBool(path)
}

func (vp *viperParser) isSet(path string) bool {
	return vp.parser.IsSet(path)
}
