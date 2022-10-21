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
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	ccjson "configcenter/src/common/json"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/BurntSushi/toml"
	"github.com/olivere/elastic/v7"
	"github.com/rwynn/monstache/monstachemap"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// blueking cmdb elastic monstache plugin.
// build: go build -buildmode=plugin -o monstache-plugin.so plugin.go

// elastic index versions.
// NOTE: CHANGE the version name if you have modify the indexes metadata struct.
const (
	indexVersionBizSet         = "20210710"
	indexVersionBiz            = "20210710"
	indexVersionSet            = "20210710"
	indexVersionModule         = "20210710"
	indexVersionHost           = "20210710"
	indexVersionModel          = "20210710"
	indexVersionObjectInstance = "20210710"
)

const (
	// default metaId
	nullMetaId      = "0"
	mongoMetaId     = "_id"
	mongoOptionId   = "id"
	mongoOptionName = "name"
	mongoEnum       = "enum"
	mongoDatabase   = "cmdb"
	configPath      = "./etc/extra.toml"
	commonObject    = "common"
)

// elastic indexes.
var (
	indexBizSet         *meta.ESIndex
	indexBiz            *meta.ESIndex
	indexSet            *meta.ESIndex
	indexModule         *meta.ESIndex
	indexHost           *meta.ESIndex
	indexModel          *meta.ESIndex
	indexObjectInstance *meta.ESIndex
	indexList           []*meta.ESIndex
)

type extraConfig struct {
	// assign es replicaNum
	ReplicaNum string `toml:"elasticsearch-shard-num"`

	// assign es shardingNum
	ShardingNum string `toml:"elasticsearch-replica-num"`
}

type instEnumIdToName struct {
	// the struct like: map[obj]map[bk_property_id]map[option.id]option.name
	instEnumMap map[string]map[string]map[string]string
	rw          sync.RWMutex
}

// 对于资源池等内部资源需要进行屏蔽，由于名字可能会改，所以需要通过cc_ApplicationBase表中的 defaultId为1进行判断
type skipBizId struct {
	bizIds map[int64]struct{}
	rw     sync.RWMutex
}

type bizId struct {
	BusinessID int64 `json:"bk_biz_id" bson:"bk_biz_id"`
}

var (
	instEnumInfo  *instEnumIdToName
	skipBizIdList *skipBizId
)

// cronInsEnumInfo TODO
// regular update instance enum ID to name
func cronInsEnumInfo(input *monstachemap.InitPluginInput) {

	// init object instance option's "id" to "name"
	initEnum := func() {
		instEnumInfoTmp := &instEnumIdToName{
			instEnumMap: make(map[string]map[string]map[string]string),
		}
		// step 1 : search all models
		models := make([]map[string]interface{}, 0)
		modelCursor, err := input.MongoClient.Database(mongoDatabase).Collection(common.BKTableNameObjDes).
			Find(context.Background(), bson.D{})
		if err != nil {
			log.Printf("query model attributes cursor failed, err: %v", err)
			return
		}

		if err := modelCursor.All(context.Background(), &models); err != nil {
			log.Printf("query model attributes failed, err: %v", err)
			return
		}

		objIds := make([]string, 0)
		for _, model := range models {
			if obj, ok := model[common.BKObjIDField].(string); ok {
				objIds = append(objIds, obj)
			}
		}

		// step 2： search  enum and bk_property_id in model attribute.
		for _, obj := range objIds {

			// query model attribute.
			modelAttrs := make([]map[string]interface{}, 0)

			modelAttrsCursor, err := input.MongoClient.Database(mongoDatabase).Collection(common.BKTableNameObjAttDes).
				Find(context.Background(), bson.D{{common.BKObjIDField, obj},
					{common.BKPropertyTypeField, mongoEnum}})
			if err != nil {
				return
			}
			if err := modelAttrsCursor.All(context.Background(), &modelAttrs); err != nil {
				return
			}
			instEnumInfoTmp.instEnumMap[obj] = make(map[string]map[string]string)
			tmpPropertyIDMap := make(map[string]map[string]string)

			for _, modelAttr := range modelAttrs {
				optionMap := make(map[string]string)
				if _, ok := modelAttr[common.BKPropertyIDField].(string); !ok {
					continue
				}
				if attr, ok := modelAttr[common.BKOptionField].(primitive.A); ok {
					opts := []interface{}(attr)
					for _, opt := range opts {
						// option.id:option.name
						if o, ok := opt.(map[string]interface{}); ok {
							if _, ok := o[mongoOptionName].(string); ok {
								optionMap[o[mongoOptionId].(string)] = o[mongoOptionName].(string)
							}
						}
					}
				}
				tmpPropertyIDMap[modelAttr[common.BKPropertyIDField].(string)] = optionMap
			}
			instEnumInfoTmp.instEnumMap[obj] = tmpPropertyIDMap
		}

		instEnumInfo.rw.Lock()
		defer instEnumInfo.rw.Unlock()
		instEnumInfo.instEnumMap = instEnumInfoTmp.instEnumMap
		log.Println("update instEnumInfo successfully")
		return
	}

	for {
		initEnum()
		err := initSkipBizId(input)
		if err != nil {
			log.Printf("init resource pool fail, err: %v", err)
			os.Exit(1)
		}
		time.Sleep(time.Minute)
	}

}

// newESIndexMetadata new es index metadata
func newESIndexMetadata(config extraConfig) *meta.ESIndexMetadata {
	return &meta.ESIndexMetadata{
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
					PropertyType: meta.IndexPropertyTypeKeyword,
				},
			},
		},
	}
}

func init() {
	// initialize each index for this release version plugin.
	var config extraConfig
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		panic(err)
	}

	if config.ShardingNum == "" || config.ReplicaNum == "" {
		panic(fmt.Sprintln("es shardingNum or replicaNum is not config!"))
	}

	instEnumInfo = &instEnumIdToName{
		instEnumMap: make(map[string]map[string]map[string]string),
	}
	skipBizIdList = &skipBizId{
		bizIds: make(map[int64]struct{}),
	}

	// biz set index.
	indexBizSetMetadata := newESIndexMetadata(config)
	indexBizSetMetadata.Mappings.Properties[meta.IndexPropertyBKBizSetID] = meta.ESIndexMetaMappingsProperty{
		PropertyType: meta.IndexPropertyTypeKeyword,
	}
	// init indexBizSetMetadata, but biz set not meta.IndexPropertyBKBizID field, delete it
	delete(indexBizSetMetadata.Mappings.Properties, meta.IndexPropertyBKBizID)
	indexBizSet = meta.NewESIndex(meta.IndexNameBizSet, indexVersionBizSet, indexBizSetMetadata)
	indexList = append(indexList, indexBizSet)

	// business application index.
	indexBizMetadata := newESIndexMetadata(config)
	indexBiz = meta.NewESIndex(meta.IndexNameBiz, indexVersionBiz, indexBizMetadata)
	indexList = append(indexList, indexBiz)

	// set index.
	indexSetMetadata := newESIndexMetadata(config)
	indexSetMetadata.Mappings.Properties[meta.IndexPropertyBKParentID] = meta.ESIndexMetaMappingsProperty{
		PropertyType: meta.IndexPropertyTypeKeyword,
	}
	indexSet = meta.NewESIndex(meta.IndexNameSet, indexVersionSet, indexSetMetadata)
	indexList = append(indexList, indexSet)

	// module index.
	indexModuleMetadata := newESIndexMetadata(config)
	indexModule = meta.NewESIndex(meta.IndexNameModule, indexVersionModule, indexModuleMetadata)
	indexList = append(indexList, indexModule)

	// host index.
	indexHostMetadata := newESIndexMetadata(config)
	indexHostMetadata.Mappings.Properties[meta.IndexPropertyBKCloudID] = meta.ESIndexMetaMappingsProperty{
		PropertyType: meta.IndexPropertyTypeKeyword,
	}
	// init indexHostMetadata, but host is not meta.IndexPropertyBKBizID field, delete it
	delete(indexHostMetadata.Mappings.Properties, meta.IndexPropertyBKBizID)
	indexHost = meta.NewESIndex(meta.IndexNameHost, indexVersionHost, indexHostMetadata)
	indexList = append(indexList, indexHost)

	// model index.
	indexModelMetadata := newESIndexMetadata(config)
	indexModel = meta.NewESIndex(meta.IndexNameModel, indexVersionModel, indexModelMetadata)
	indexList = append(indexList, indexModel)

	// object instance index.
	indexObjInstMetadata := newESIndexMetadata(config)
	indexObjectInstance = meta.NewESIndex(meta.IndexNameObjectInstance, indexVersionObjectInstance, indexObjInstMetadata)
	indexList = append(indexList, indexObjectInstance)

	log.Println("bk-cmdb elastic monstache plugin initialize successfully")
}

// analysisJSONKeywords analysis the given json style document, and extract
// all the keywords as elastic document content.
func analysisJSONKeywords(result gjson.Result) []string {

	keywords := make([]string, 0)
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

	compressedKeywords := make([]string, 0)
	// keywordsMap control repeated or screened keywords.
	keywordsMap := make(map[string]struct{})
	for _, keyword := range keywords {
		if _, exist := keywordsMap[keyword]; exist {
			continue
		}
		compressedKeywords = append(compressedKeywords, keyword)
		keywordsMap[keyword] = struct{}{}
	}

	return compressedKeywords
}

// getMetaIdToStr objID/hostID/setID/moduleID/instanceID/bizID  convert to string.
func getMetaIdToStr(d interface{}) (string, error) {
	if d == nil {
		return "", errors.New("document id is nil")
	}
	return fmt.Sprintf("%v", d), nil
}

// baseDataCleaning  do not need to sync "_id","create_time","last_time","bk_supplier_account".
func baseDataCleaning(document map[string]interface{}) map[string]interface{} {
	delete(document, mongoMetaId)
	delete(document, common.CreateTimeField)
	delete(document, common.LastTimeField)
	delete(document, common.BKOwnerIDField)
	return document
}

// originalDataCleaning some field do not need to save es,delete it.
func originalDataCleaning(document map[string]interface{}, collection string) map[string]interface{} {

	if document == nil {
		return nil
	}

	doc := make(map[string]interface{})

	switch collection {
	case common.BKTableNameBaseBizSet:
		doc = baseDataCleaning(document)
		// do not need to sync "default".
		delete(doc, common.BKDefaultField)
		delete(doc, common.BKBizSetScopeField)

	case common.BKTableNameBaseApp:
		doc = baseDataCleaning(document)
		// do not need to sync "default".
		delete(doc, common.BKDefaultField)
		delete(doc, common.BKParentIDField)

	case common.BKTableNameBaseSet:

		doc = baseDataCleaning(document)

		// do not need to sync "default","set_template_id","bk_biz_id","bk_parent_id".
		delete(doc, common.BKAppIDField)
		delete(doc, common.BKParentIDField)
		delete(doc, common.BKSetTemplateIDField)
		delete(doc, common.BKDefaultField)

	case common.BKTableNameBaseModule:
		doc = baseDataCleaning(document)

		// do not need to sync "default","set_template_id","bk_biz_id","bk_parent_id","bk_set_id","service_category_id".
		delete(doc, common.BKDefaultField)
		delete(doc, common.BKSetTemplateIDField)
		delete(doc, common.BKAppIDField)
		delete(doc, common.BKParentIDField)
		delete(doc, common.BKSetIDField)
		delete(doc, common.BKServiceCategoryIDField)

	case common.BKTableNameBaseHost:

		doc = baseDataCleaning(document)
		// do not need to sync "operation_time".
		delete(doc, common.BKOperationTimeField)
		delete(doc, common.BKParentIDField)

	case common.BKTableNameObjDes:

		// need to sync "bk_obj_name" and "bk_obj_id".
		doc[common.BKObjIDField] = document[common.BKObjIDField]
		doc[common.BKObjNameField] = document[common.BKObjNameField]

	case common.BKTableNameObjAttDes:

		// need to sync "bk_property_id" and "bk_property_name".
		doc[common.BKPropertyIDField] = document[common.BKPropertyIDField]
		doc[common.BKPropertyNameField] = document[common.BKPropertyNameField]

	default:
		doc = baseDataCleaning(document)
		// do not need to sync "bk_obj_id" for common object instance.
		delete(doc, common.BKObjIDField)
		delete(doc, common.BKParentIDField)
	}

	return doc
}

// getModeNameByCollection parse the innerObjId  from collection name.
func getModeNameByCollection(collection string) (innerObjId string) {

	switch collection {
	case common.BKTableNameBaseBizSet:
		innerObjId = common.BKInnerObjIDBizSet
	case common.BKTableNameBaseHost:
		innerObjId = common.BKInnerObjIDHost
	case common.BKTableNameBaseApp:
		innerObjId = common.BKInnerObjIDApp
	case common.BKTableNameBaseSet:
		innerObjId = common.BKInnerObjIDSet

	case common.BKTableNameBaseModule:
		innerObjId = common.BKInnerObjIDModule
	default:
		if common.IsObjectInstShardingTable(collection) {
			tmp := strings.TrimLeft(collection, common.BKObjectInstShardingTablePrefix)
			instSlice := strings.Split(tmp, "_")
			if len(instSlice) >= 3 {
				innerObjId = strings.Join(instSlice[2:], "_")
			}
		}
	}
	return innerObjId
}

// enumIdToName parse enum Id to Name.
func enumIdToName(document map[string]interface{}, collection string) {

	key := getModeNameByCollection(collection)
	if key == "" {
		return
	}
	instEnumInfo.rw.RLock()
	defer instEnumInfo.rw.RUnlock()
	// deal enum  map[string]map[string]map[string]string
	for propertyId, enumInfo := range instEnumInfo.instEnumMap[key] {
		if _, ok := document[propertyId]; ok {
			if v, ok := document[propertyId].(string); ok {
				document[propertyId] = enumInfo[v]
			}
		}
	}
	return
}

// analysisDocument analysis the given document, return document id and keywords.
func analysisDocument(document map[string]interface{}, collection string) (string, []string, error) {

	var id string
	// analysis collection document id.
	switch collection {
	case common.BKTableNameBaseBizSet:
		bizSetId, err := getMetaIdToStr(document[common.BKBizSetIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKBizSetIDField, err)
		}
		id = bizSetId
	case common.BKTableNameBaseApp:
		bizId, err := getMetaIdToStr(document[common.BKAppIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKAppIDField, err)
		}
		id = bizId
	case common.BKTableNameBaseSet:

		setId, err := getMetaIdToStr(document[common.BKSetIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKSetIDField, err)
		}
		id = setId
	case common.BKTableNameBaseModule:

		moduleId, err := getMetaIdToStr(document[common.BKModuleIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKModuleIDField, err)
		}
		id = moduleId

	case common.BKTableNameBaseHost:
		hostId, err := getMetaIdToStr(document[common.BKHostIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKHostIDField, err)
		}
		id = hostId

	case common.BKTableNameObjDes, common.BKTableNameObjAttDes:
		objId, err := getMetaIdToStr(document[common.BKObjIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKObjIDField, err)
		}
		id = objId
	default:
		instId, err := getMetaIdToStr(document[common.BKInstIDField])
		if err != nil {
			return "", nil, fmt.Errorf("missing: %s, err: %v", common.BKInstIDField, err)
		}
		id = instId
	}
	// in the instance scenario, the enumeration values need to be converted
	if collection != common.BKTableNameObjDes {
		enumIdToName(document, collection)
	}

	doc := originalDataCleaning(document, collection)
	if doc == nil {
		return "", nil, errors.New("there is no document")
	}
	// analysis keywords.
	jsonDoc, err := ccjson.MarshalToString(doc)
	if err != nil {
		return "", nil, err
	}
	keywords := analysisJSONKeywords(gjson.Parse(jsonDoc))

	// return document id and compressed keywords.
	return id, compressKeywords(keywords), nil
}

// outputDocument return output document
func outputDocument(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput, objID,
	esObjID string) (map[string]interface{}, error) {
	oId := input.Document[common.BKOwnerIDField]
	metaId := input.Document[mongoMetaId]
	bizId := input.Document[common.BKAppIDField]

	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return nil, fmt.Errorf("analysis output document failed, document: %+v, err: %v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           objID,
		meta.IndexPropertyBKSupplierAccount: oId,
		meta.IndexPropertyBKBizID:           bizId,
		meta.IndexPropertyKeywords:          keywords,
	}

	documentID, ok := metaId.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("missing document metadata id")
	}
	idEs := fmt.Sprintf("%s:%s", documentID.Hex(), esObjID)
	output.ID = idEs

	return document, nil
}

// indexingBizSet indexing the business set instance.
func indexingBizSet(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	bizSetId := input.Document[common.BKBizSetIDField]

	document, err := outputDocument(input, output, common.BKInnerObjIDBizSet, common.BKInnerObjIDBizSet)
	if err != nil {
		return fmt.Errorf("get biz set output document failed, err: %v", err)
	}
	document[meta.IndexPropertyBKBizSetID] = bizSetId
	delete(document, meta.IndexPropertyBKBizID)

	output.Document = document
	// use alias name to indexing document.
	output.Index = indexBizSet.AliasName()

	return nil
}

// indexingApplication indexing the business application instance.
func indexingApplication(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	document, err := outputDocument(input, output, common.BKInnerObjIDApp, common.BKInnerObjIDApp)
	if err != nil {
		return fmt.Errorf("get biz output document failed, err: %v", err)
	}

	output.Document = document
	// use alias name to indexing document.
	output.Index = indexBiz.AliasName()

	return nil
}

// indexingSet indexing the set instance.
func indexingSet(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	pId := input.Document[common.BKParentIDField]

	document, err := outputDocument(input, output, common.BKInnerObjIDSet, common.BKInnerObjIDSet)
	if err != nil {
		return fmt.Errorf("get set output document failed, err: %v", err)
	}
	document[meta.IndexPropertyBKParentID] = pId

	output.Document = document
	// use alias name to indexing document.
	output.Index = indexSet.AliasName()

	return nil
}

// indexingModule indexing the module instance.
func indexingModule(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	document, err := outputDocument(input, output, common.BKInnerObjIDModule, common.BKInnerObjIDModule)
	if err != nil {
		return fmt.Errorf("get module output document failed, err: %v", err)
	}

	output.Document = document
	// use alias name to indexing document.
	output.Index = indexModule.AliasName()

	return nil
}

// indexingHost indexing the host instance.
func indexingHost(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	document, err := outputDocument(input, output, common.BKInnerObjIDHost, common.BKInnerObjIDHost)
	if err != nil {
		return fmt.Errorf("get host output document failed, err: %v", err)
	}
	document[meta.IndexPropertyBKCloudID] = input.Document[common.BKCloudIDField]
	delete(document, meta.IndexPropertyBKBizID)

	output.Document = document
	// use alias name to indexing document.
	output.Index = indexHost.AliasName()

	return nil
}

// indexingModel indexing the model/attr instance.
func indexingModel(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	objectID, ok := input.Document[common.BKObjIDField].(string)
	if !ok {
		return fmt.Errorf("analysis model document failed, object id missing, %+v", input.Document)
	}

	// query model.
	model := make(map[string]interface{})

	if err := input.MongoClient.Database(input.Database).Collection(common.BKTableNameObjDes).
		FindOne(context.Background(), bson.D{{common.BKObjIDField, objectID}}).
		Decode(&model); err != nil {
		return fmt.Errorf("query model object[%s] failed, %v", objectID, err)
	}

	oId := model[common.BKOwnerIDField]
	bizId := model[common.BKAppIDField]
	metaId := model[mongoMetaId]

	// analysis model document.
	_, keywords, err := analysisDocument(model, common.BKTableNameObjDes)
	if err != nil {
		return fmt.Errorf("analysis model document failed, %+v, %v", input.Document, err)
	}

	// query model attribute.
	modelAttrs := make([]map[string]interface{}, 0)

	modelAttrsCursor, err := input.MongoClient.Database(input.Database).Collection(common.BKTableNameObjAttDes).
		Find(context.Background(), bson.D{{common.BKObjIDField, objectID}})
	if err != nil {
		return fmt.Errorf("query model attributes object[%s] cursor failed, %v", objectID, err)
	}

	if err := modelAttrsCursor.All(context.Background(), &modelAttrs); err != nil {
		return fmt.Errorf("query model attributes object[%s] failed, %v", objectID, err)
	}

	// all attributes with model metadata is ONE elastic document.
	for _, attribute := range modelAttrs {
		// data Cleaning
		attr := originalDataCleaning(attribute, common.BKTableNameObjAttDes)
		jsonDoc, err := ccjson.MarshalToString(attr)
		if err != nil {
			log.Printf("marshal model attributes object[%s] failed, %+v, %v", objectID, attribute, err)
			continue
		}
		keywords = append(keywords, analysisJSONKeywords(gjson.Parse(jsonDoc))...)
	}
	documentID, ok := metaId.(primitive.ObjectID)
	if !ok {
		return errors.New("missing document metadata id")
	}
	idEs := fmt.Sprintf("%s:%s", documentID.Hex(), common.BKInnerObjIDObject)

	// build elastic document.
	document := map[string]interface{}{
		// model scene,we use meta_bk_obj_id to search mongo,this id set null.
		meta.IndexPropertyID:                nullMetaId,
		meta.IndexPropertyDataKind:          meta.DataKindModel,
		meta.IndexPropertyBKObjID:           objectID,
		meta.IndexPropertyBKSupplierAccount: oId,
		meta.IndexPropertyBKBizID:           bizId,
		meta.IndexPropertyKeywords:          compressKeywords(keywords),
	}
	output.ID = idEs
	output.Document = document
	// use alias name to indexing document.
	output.Index = indexModel.AliasName()

	return nil
}

// indexingObjectInstance indexing the common object instance.
func indexingObjectInstance(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {

	objId := input.Document[common.BKObjIDField]
	bizId := input.Document[common.BKAppIDField]
	oId := input.Document[common.BKOwnerIDField]
	metaId := input.Document[mongoMetaId]

	// analysis document.
	id, keywords, err := analysisDocument(input.Document, input.Collection)
	if err != nil {
		return fmt.Errorf("analysis object instance document failed, %+v, %v", input.Document, err)
	}

	// build elastic document.
	document := map[string]interface{}{
		meta.IndexPropertyID:                id,
		meta.IndexPropertyDataKind:          meta.DataKindInstance,
		meta.IndexPropertyBKObjID:           objId,
		meta.IndexPropertyBKSupplierAccount: oId,
		meta.IndexPropertyBKBizID:           bizId,
		meta.IndexPropertyKeywords:          keywords,
	}

	documentID, ok := metaId.(primitive.ObjectID)
	if !ok {
		return errors.New("missing document metadata id")
	}
	idEs := fmt.Sprintf("%s:%s", documentID.Hex(), commonObject)
	output.ID = idEs

	output.Document = document
	// use alias name to indexing document.
	output.Index = indexObjectInstance.AliasName()

	return nil
}

// initSkipBizId TODO
// the internal resource pool does not need to be displayed externally. The ID corresponding to the internal resource
// pool is saved. When writing to es from Mongo, the relevant doc needs to be masked.
func initSkipBizId(input *monstachemap.InitPluginInput) error {

	bizInfo := make([]bizId, 0)
	appCursor, err := input.MongoClient.Database(mongoDatabase).Collection(common.BKTableNameBaseApp).
		Find(context.Background(), bson.D{{common.BKDefaultField, 1}})
	if err != nil {
		return fmt.Errorf("query app database appCursor fail, err: %v", err)
	}

	if err := appCursor.All(context.Background(), &bizInfo); err != nil {
		return fmt.Errorf("query app database fail, err: %v", err)
	}
	if len(bizInfo) == 0 {
		return errors.New("query list num is zero")
	}
	skipBizIdList.rw.Lock()
	defer skipBizIdList.rw.Unlock()

	for _, v := range bizInfo {
		skipBizIdList.bizIds[v.BusinessID] = struct{}{}
	}
	log.Printf("initSkipBizId success, bizId: %v", bizInfo)
	return nil
}

// Init function, when you implement a Init function, it would load and call this function with the initialized
// mongo/elastic clients. And you could do some initialization for elasticsearch or mongodb here.
func Init(input *monstachemap.InitPluginInput) error {

	go cronInsEnumInfo(input)

	// initialize elastic indexes.
	for _, index := range indexList {
		// check elastic index.
		exist, err := input.ElasticClient.IndexExists(index.Name()).Do(context.Background())
		if err != nil {
			return fmt.Errorf("check elastic index[%s] existence failed, %v", index.Name(), err)
		}

		if !exist {
			// NOTE: create new index with the target index name, and it may be a alias index name,
			// the policies are all by user.
			_, err = input.ElasticClient.CreateIndex(index.Name()).Body(index.Metadata()).Do(context.Background())
			if err != nil {
				return fmt.Errorf("create elastic index[%s] failed, %v", index.Name(), err)
			}
		}

		// check elastic alias name.
		// it's ok if the alias name index is already exist, but the alias name could not be a real index.
		_, err = input.ElasticClient.Alias().Add(index.Name(), index.AliasName()).Do(context.Background())
		if err != nil {
			return fmt.Errorf("create elastic index[%s] alias failed, %v", index.Name(), err)
		}
	}

	log.Println("initialize elastic indexes successfully")

	return nil
}

// Map function, when you implement a Map function, you could handle each event document base on the
// plugin input, the input parameter will contain information about the document's origin database and
// collection, and mapping the elastic index document in output.
func Map(input *monstachemap.MapperPluginInput) (*monstachemap.MapperPluginOutput, error) {

	defer func() {
		if errRecover := recover(); errRecover != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			log.Printf("map data panic, buf: %v", string(buf))
		}
	}()

	// discard all internal resource pool class docs.
	if input.Collection == common.BKTableNameBaseApp || input.Collection == common.BKTableNameBaseSet {
		bizId := input.Document[common.BKAppIDField]
		if bizId != nil {
			skipBizIdList.rw.RLock()
			defer skipBizIdList.rw.RUnlock()
			bId, err := util.GetInt64ByInterface(bizId)
			if err != nil {
				log.Printf("bizId convert fail, bizId: %v, err: %v", bizId, err)
				return nil, err
			}

			if _, exist := skipBizIdList.bizIds[bId]; exist {
				return nil, nil
			}
		}
	}

	output := new(monstachemap.MapperPluginOutput)
	switch input.Collection {
	case common.BKTableNameBaseBizSet:
		if err := indexingBizSet(input, output); err != nil {
			return nil, err
		}

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

// Process function, when you implement a Process function, the function will be called after monstache processes each
// event. This function has full access to the MongoDB and Elasticsearch clients (
// including the Elasticsearch bulk processor) in the input and allows you to handle complex event processing scenarios
func Process(input *monstachemap.ProcessPluginInput) error {
	req := elastic.NewBulkDeleteRequest()
	metaId := input.Document[mongoMetaId]
	documentID, ok := metaId.(primitive.ObjectID)
	if !ok {
		return errors.New("missing document metadata id")
	}

	var index, objectID string
	objectID = documentID.Hex()
	if input.Operation == "d" {
		switch input.Collection {
		case common.BKTableNameBaseBizSet:
			objectID = fmt.Sprintf("%s:%s", objectID, common.BKInnerObjIDBizSet)
			index = fmt.Sprintf("%s_%s", meta.IndexNameBizSet, indexVersionBizSet)
		case common.BKTableNameBaseApp:
			objectID = fmt.Sprintf("%s:%s", objectID, common.BKInnerObjIDApp)
			index = fmt.Sprintf("%s_%s", meta.IndexNameBiz, indexVersionBiz)
		case common.BKTableNameBaseSet:
			objectID = fmt.Sprintf("%s:%s", objectID, common.BKInnerObjIDSet)
			index = fmt.Sprintf("%s_%s", meta.IndexNameSet, indexVersionSet)
		case common.BKTableNameBaseModule:
			objectID = fmt.Sprintf("%s:%s", objectID, common.BKInnerObjIDModule)
			index = fmt.Sprintf("%s_%s", meta.IndexNameModule, indexVersionModule)
		case common.BKTableNameBaseHost:
			objectID = fmt.Sprintf("%s:%s", objectID, common.BKInnerObjIDHost)
			index = fmt.Sprintf("%s_%s", meta.IndexNameHost, indexVersionHost)
		case common.BKTableNameObjDes, common.BKTableNameObjAttDes:
			objectID = fmt.Sprintf("%s:%s", objectID, common.BKInnerObjIDObject)
			index = fmt.Sprintf("%s_%s", meta.IndexNameModel, indexVersionModel)
		default:
			objectID = fmt.Sprintf("%s:%s", objectID, commonObject)
			index = fmt.Sprintf("%s_%s", meta.IndexNameObjectInstance, indexVersionObjectInstance)
		}

		req.Id(objectID)
		req.Index(index)
		input.ElasticBulkProcessor.Add(req)
	}

	return nil
}
