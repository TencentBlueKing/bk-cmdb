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

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"configcenter/src/common"
	ccjson "configcenter/src/common/json"
	meta "configcenter/src/common/metadata"

	"github.com/BurntSushi/toml"
	"github.com/rwynn/monstache/monstachemap"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

// blueking cmdb elastic monstache plugin.
// build: go build -buildmode=plugin -o monstache-plugin.so plugin.go

// elastic index versions.
// NOTE: CHANGE the version name if you have modify the indexes metadata struct.
const (
	indexVersionBiz            = "20210710"
	indexVersionSet            = "20210710"
	indexVersionModule         = "20210710"
	indexVersionHost           = "20210710"
	indexVersionModel          = "20210710"
	indexVersionObjectInstance = "20210710"
)

const (
	// default metaId
	nullMetaId     = "0"
	mongoMetaId    = "_id"
	mongoCreatTime = "create_time"
	mongoLastTime  = "last_time"
	configPath     = "./etc/extra.toml"
)

// elastic indexes.
var (
	indexBiz            *meta.ESIndex
	indexSet            *meta.ESIndex
	indexModule         *meta.ESIndex
	indexHost           *meta.ESIndex
	indexModel          *meta.ESIndex
	indexObjectInstance *meta.ESIndex
	indexList           []*meta.ESIndex
)

type extraConfig struct {
	ReplicaNum  string `toml:"elasticsearch-shard-num"`
	ShardingNum string `toml:"elasticsearch-replica-num"`
}

func init() {
	// initialize each index for this release version plugin.
	var config extraConfig
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		panic(err)
	}
	if config.ShardingNum == "" || config.ReplicaNum == "" {
		panic(fmt.Sprintf("es shardingNum or replicaNum is not config!"))
	}

	// business application index.
	indexBiz = meta.NewESIndex(meta.IndexNameBiz, indexVersionBiz, &meta.ESIndexMetadata{
		Settings: meta.ESIndexMetaSettings{
			Shards:   config.ShardingNum,
			Replicas: config.ReplicaNum,
		},
		Mappings: meta.ESIndexMetaMappings{
			Properties: map[string]meta.ESIndexMetaMappingsProperty{
				meta.IndexPropertyID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKObjID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKSupplierAccount: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKBizID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyKeywords: {
					PropertyType: meta.IndexPropertyTypeText,
				},
			},
		},
	})
	indexList = append(indexList, indexBiz)

	// set index.
	indexSet = meta.NewESIndex(meta.IndexNameSet, indexVersionSet, &meta.ESIndexMetadata{
		Settings: meta.ESIndexMetaSettings{
			Shards:   config.ShardingNum,
			Replicas: config.ReplicaNum,
		},
		Mappings: meta.ESIndexMetaMappings{
			Properties: map[string]meta.ESIndexMetaMappingsProperty{
				meta.IndexPropertyID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKObjID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKSupplierAccount: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKBizID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKParentID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyKeywords: {
					PropertyType: meta.IndexPropertyTypeText,
				},
			},
		},
	})
	indexList = append(indexList, indexSet)

	// module index.
	indexModule = meta.NewESIndex(meta.IndexNameModule, indexVersionModule, &meta.ESIndexMetadata{
		Settings: meta.ESIndexMetaSettings{
			Shards:   config.ShardingNum,
			Replicas: config.ReplicaNum,
		},
		Mappings: meta.ESIndexMetaMappings{
			Properties: map[string]meta.ESIndexMetaMappingsProperty{
				meta.IndexPropertyID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKObjID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKSupplierAccount: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKBizID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyKeywords: {
					PropertyType: meta.IndexPropertyTypeText,
				},
			},
		},
	})
	indexList = append(indexList, indexModule)

	// host index.
	indexHost = meta.NewESIndex(meta.IndexNameHost, indexVersionHost, &meta.ESIndexMetadata{
		Settings: meta.ESIndexMetaSettings{
			Shards:   config.ShardingNum,
			Replicas: config.ReplicaNum,
		},
		Mappings: meta.ESIndexMetaMappings{
			Properties: map[string]meta.ESIndexMetaMappingsProperty{
				meta.IndexPropertyID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKObjID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKSupplierAccount: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKCloudID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyKeywords: {
					PropertyType: meta.IndexPropertyTypeText,
				},
			},
		},
	})
	indexList = append(indexList, indexHost)

	// model index.
	indexModel = meta.NewESIndex(meta.IndexNameModel, indexVersionModel, &meta.ESIndexMetadata{
		Settings: meta.ESIndexMetaSettings{
			Shards:   config.ShardingNum,
			Replicas: config.ReplicaNum,
		},
		Mappings: meta.ESIndexMetaMappings{
			Properties: map[string]meta.ESIndexMetaMappingsProperty{
				meta.IndexPropertyID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKObjID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKSupplierAccount: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKBizID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyKeywords: {
					PropertyType: meta.IndexPropertyTypeText,
				},
			},
		},
	})
	indexList = append(indexList, indexModel)

	// object instance index.
	indexObjectInstance = meta.NewESIndex(meta.IndexNameObjectInstance, indexVersionObjectInstance, &meta.ESIndexMetadata{
		Settings: meta.ESIndexMetaSettings{
			Shards:   config.ShardingNum,
			Replicas: config.ReplicaNum,
		},
		Mappings: meta.ESIndexMetaMappings{
			Properties: map[string]meta.ESIndexMetaMappingsProperty{
				meta.IndexPropertyID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKObjID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKSupplierAccount: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyBKBizID: {
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
				meta.IndexPropertyKeywords: {
					PropertyType: meta.IndexPropertyTypeText,
				},
			},
		},
	})
	indexList = append(indexList, indexObjectInstance)

	log.Println("bk-cmdb elastic monstache plugin initialize successfully")
}

// analysisJSONKeywords analysis the given json style document, and extract
// all the keywords as elastic document content.
func analysisJSONKeywords(result gjson.Result) []string {
	var keywords []string

	if !result.IsObject() && !result.IsArray() {
		keywords = append(keywords, result.String())
		return keywords
	}

	result.ForEach(func(key, value gjson.Result) bool {
		keywords = append(keywords, analysisJSONKeywords(value)...)
		return true
	})

	return keywords
}

// compressKeywords compress the keywords return without repetition.
func compressKeywords(keywords []string) []string {
	var compressedKeywords []string

	// keywordsMap control repeated or screened keywords.
	keywordsMap := map[string]struct{}{
		// filter empty keyword.
		"": struct{}{},
	}

	for _, keyword := range keywords {
		if _, exist := keywordsMap[keyword]; exist {
			continue
		}
		compressedKeywords = append(compressedKeywords, keyword)
		keywordsMap[keyword] = struct{}{}
	}

	return compressedKeywords
}

func getMetaIdToStr(d interface{}) (string, error) {
	if d == nil {
		return "", errors.New(fmt.Sprintf("document id is nil "))
	}
	var (
		s   string
		err error
	)
	v := reflect.ValueOf(d)
	t := v.Type()

	switch t.Kind() {
	case reflect.Int:
		s = strconv.Itoa(int(d.(int)))
	case reflect.Int16:
		s = strconv.Itoa(int(d.(int16)))
	case reflect.Int32:
		s = strconv.Itoa(int(d.(int32)))
	case reflect.Int64:
		s = strconv.FormatInt(d.(int64), 10)
	case reflect.Uint:
		s = strconv.FormatUint(uint64(d.(uint)), 10)
	case reflect.Uint32:
		s = strconv.FormatUint(uint64(d.(uint32)), 10)
	case reflect.Uint64:
		s = strconv.FormatUint(d.(uint64), 10)
	case reflect.Float32:
		s = strconv.FormatFloat(float64(d.(float32)), 'E', -1, 32)
	case reflect.Float64:
		s = strconv.FormatFloat(d.(float64), 'E', -1, 64)
	case reflect.String:
		s = d.(string)
	default:
		err = errors.New(fmt.Sprintf("type is error,type: %v", t))
	}
	return s, err
}

// analysisDocument analysis the given document, return document id and keywords.
func analysisDocument(document map[string]interface{}, colletion string) (string, []string, error) {

	var id string

	// analysis collection document id.
	switch colletion {
	case common.BKTableNameBaseApp:

		bizId, err := getMetaIdToStr(document[common.BKAppIDField])
		if err != nil {
			return "", nil, errors.New(fmt.Sprintf("missing: %s,err: %v", common.BKAppIDField, err))
		}
		id = bizId

	case common.BKTableNameBaseSet:

		setId, err := getMetaIdToStr(document[common.BKSetIDField])
		if err != nil {
			return "", nil, errors.New(fmt.Sprintf("missing: %s,err: %v", common.BKSetIDField, err))
		}
		id = setId
	case common.BKTableNameBaseModule:

		moduleId, err := getMetaIdToStr(document[common.BKModuleIDField])
		if err != nil {
			return "", nil, errors.New(fmt.Sprintf("missing: %s,err: %v", common.BKModuleIDField, err))
		}
		id = moduleId

	case common.BKTableNameBaseHost:
		hostId, err := getMetaIdToStr(document[common.BKHostIDField])
		if err != nil {
			return "", nil, errors.New(fmt.Sprintf("missing: %s,err: %v", common.BKHostIDField, err))
		}
		id = hostId
	case common.BKTableNameObjDes, common.BKTableNameObjAttDes:

		objId, err := getMetaIdToStr(document[common.BKObjIDField])
		if err != nil {
			return "", nil, errors.New(fmt.Sprintf("missing: %s,err: %v", common.BKObjIDField, err))
		}
		id = objId
	default:
		instId, err := getMetaIdToStr(document[common.BKInstIDField])
		if err != nil {

			return "", nil, errors.New(fmt.Sprintf("missing: %s,err: %v", common.BKInstIDField, err))
		}
		id = instId
	}

	// metaID creattime lasttime no need to store
	if document[mongoMetaId] != nil {
		delete(document, mongoMetaId)
	}
	if document[mongoCreatTime] != nil {
		delete(document, mongoCreatTime)
	}
	if document[mongoLastTime] != nil {
		delete(document, mongoLastTime)
	}

	// analysis keywords.
	jsonDoc, err := ccjson.MarshalToString(document)
	if err != nil {
		return "", nil, err
	}
	keywords := analysisJSONKeywords(gjson.Parse(jsonDoc))

	// return document id and compressed keywords.
	return id, compressKeywords(keywords), nil
}

// indexingApplication indexing the business application instance.
func indexingApplication(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis business application document failed, %+v, %+v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           common.BKInnerObjIDApp,
		meta.IndexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		meta.IndexPropertyBKBizID:           input.Document[common.BKAppIDField],
		meta.IndexPropertyKeywords:          keywords,
	}

	output.ID = id
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexBiz.AliasName()

	return nil
}

// indexingSet indexing the set instance.
func indexingSet(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis set document failed, %+v, %+v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           common.BKInnerObjIDSet,
		meta.IndexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		meta.IndexPropertyBKBizID:           input.Document[common.BKAppIDField],
		meta.IndexPropertyBKParentID:        input.Document[common.BKParentIDField],
		meta.IndexPropertyKeywords:          keywords,
	}

	output.ID = id
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexSet.AliasName()

	return nil
}

// indexingModule indexing the module instance.
func indexingModule(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis module document failed, %+v, %+v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           common.BKInnerObjIDModule,
		meta.IndexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		meta.IndexPropertyBKBizID:           input.Document[common.BKAppIDField],
		meta.IndexPropertyKeywords:          keywords,
	}

	output.ID = id
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexModule.AliasName()

	return nil
}

// indexingHost indexing the host instance.
func indexingHost(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis host document failed, %+v, %+v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           common.BKInnerObjIDHost,
		meta.IndexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		meta.IndexPropertyBKCloudID:         input.Document[common.BKCloudIDField],
		meta.IndexPropertyKeywords:          keywords,
	}

	output.ID = id
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexHost.AliasName()

	return nil
}

// indexingModel indexing the model/attr instance.
func indexingModel(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// model object id.
	objectID, ok := input.Document[common.BKObjIDField].(string)
	if !ok {
		return fmt.Errorf("analysis model document failed, object id missing, %+v", input.Document)
	}

	// query model.
	model := make(map[string]interface{})

	if err := input.MongoClient.Database(input.Database).Collection(common.BKTableNameObjDes).
		FindOne(context.Background(), bson.D{{common.BKObjIDField, objectID}}).
		Decode(&model); err != nil {
		return fmt.Errorf("query model object[%s] failed, %+v", objectID, err)
	}

	// analysis document.
	id, keywords, err := analysisDocument(model, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis model document failed, %+v, %+v", input.Document, err)
	}

	// query model attribute.
	modelAttrs := []map[string]interface{}{}

	modelAttrsCursor, err := input.MongoClient.Database(input.Database).Collection(common.BKTableNameObjAttDes).
		Find(context.Background(), bson.D{{common.BKObjIDField, objectID}})
	if err != nil {
		return fmt.Errorf("query model attributes object[%s] cursor failed, %+v", objectID, err)
	}

	if err := modelAttrsCursor.All(context.Background(), &modelAttrs); err != nil {
		return fmt.Errorf("query model attributes object[%s] failed, %+v", objectID, err)
	}

	// merge model attribute keywords,
	// all attributes with model metadata is ONE elastic document.
	for _, attribute := range modelAttrs {
		jsonDoc, err := ccjson.MarshalToString(attribute)
		if err != nil {
			return fmt.Errorf("marshal model attributes object[%s] failed, %+v, %+v", objectID, attribute, err)
		}
		keywords = append(keywords, analysisJSONKeywords(gjson.Parse(jsonDoc))...)
	}

	// build elastic document.
	document := map[string]interface{}{
		// model scene,we use meta_bk_obj_id to search mongo,this id set null
		meta.IndexPropertyID:                nullMetaId,
		meta.IndexPropertyDataKind:          meta.DataKindModel,
		meta.IndexPropertyBKObjID:           objectID,
		meta.IndexPropertyBKSupplierAccount: model[common.BKOwnerIDField],
		meta.IndexPropertyBKBizID:           model[common.BKAppIDField],
		meta.IndexPropertyKeywords:          compressKeywords(keywords),
	}

	output.ID = id
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexModel.AliasName()

	return nil
}

// indexingObjectInstance indexing the common object instance.
func indexingObjectInstance(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis object instance document failed, %+v, %+v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           input.Document[common.BKObjIDField],
		meta.IndexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		meta.IndexPropertyBKBizID:           input.Document[common.BKAppIDField],
		meta.IndexPropertyKeywords:          keywords,
	}

	output.ID = id
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexObjectInstance.AliasName()

	return nil
}

// Init function, when you implement a Init function, it would load and call this function with the initialized
// mongo/elastic clients. And you could do some initialization for elasticsearch or mongodb here.
func Init(input *monstachemap.InitPluginInput) error {
	// initialize elastic indexes.
	for _, index := range indexList {
		// check elastic index.
		exist, err := input.ElasticClient.IndexExists(index.Name()).Do(context.Background())
		if err != nil {
			return fmt.Errorf("check elastic index[%s] existence failed, %+v", index.Name(), err)
		}

		if !exist {
			// NOTE: create new index with the target index name, and it may be a alias index name,
			// the policies are all by user.
			_, err = input.ElasticClient.CreateIndex(index.Name()).Body(index.Metadata()).Do(context.Background())
			if err != nil {
				return fmt.Errorf("create elastic index[%s] failed, %+v", index.Name(), err)
			}
		}

		// check elastic alias name.
		// it's ok if the alias name index is already exist, but the alias name could not be a real index.
		_, err = input.ElasticClient.Alias().Add(index.Name(), index.AliasName()).Do(context.Background())
		if err != nil {
			return fmt.Errorf("create elastic index[%s] alias failed, %+v", index.Name(), err)
		}
	}

	log.Printf("initialize elastic indexes successfully")

	return nil
}

// Map function, when you implement a Map function, you could handle each event document base on the
// plugin input, the input parameter will contain information about the document's origin database and
// collection, and mapping the elastic index document in output.
func Map(input *monstachemap.MapperPluginInput) (*monstachemap.MapperPluginOutput, error) {
	output := &monstachemap.MapperPluginOutput{}

	switch input.Collection {
	case common.BKTableNameBaseApp:
		if err := indexingApplication(input, output); err != nil {
			return nil, err
		}

	case common.BKTableNameBaseSet:
		if err := indexingSet(input, output); err != nil {
			return nil, err
		}

	case common.BKTableNameBaseModule:
		if err := indexingModule(input, output); err != nil {
			return nil, err
		}

	case common.BKTableNameBaseHost:
		if err := indexingHost(input, output); err != nil {
			return nil, err
		}

	case common.BKTableNameObjDes, common.BKTableNameObjAttDes:
		if err := indexingModel(input, output); err != nil {
			return nil, err
		}

	default:
		if !common.IsObjectShardingTable(input.Collection) {
			// unknown collection, just drop it.
			output.Drop = true
			return output, nil
		}

		if err := indexingObjectInstance(input, output); err != nil {
			return nil, err
		}
	}

	return output, nil
}
