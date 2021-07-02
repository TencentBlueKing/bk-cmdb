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

	"github.com/rwynn/monstache/monstachemap"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"

	"configcenter/src/common"
	ccjson "configcenter/src/common/json"
	"configcenter/src/common/metadata"
)

// blueking cmdb elastic monstache plugin.
// build: go build -buildmode=plugin -o bk-cmdb-monstache-plugin.so plugin.go

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

// bk-cmdb database.
const (
	database = "cmdb"
)

// elastic index names.
const (
	indexNameBiz            = "bk_cmdb.biz"
	indexNameSet            = "bk_cmdb.set"
	indexNameModule         = "bk_cmdb.module"
	indexNameHost           = "bk_cmdb.host"
	indexNameModel          = "bk_cmdb.model"
	indexNameObjectInstance = "bk_cmdb.object_instance"
)

// elastic index property types.
const (
	indexPropertyTypeKeyword = "keyword"
	indexPropertyTypeText    = "text"
)

// elastic index properties.
const (
	indexPropertyID                = "meta_id"
	indexPropertyBKObjID           = "meta_bk_obj_id"
	indexPropertyBKSupplierAccount = "meta_bk_supplier_account"
	indexPropertyBKBizID           = "meta_bk_biz_id"
	indexPropertyBKParentID        = "meta_bk_parent_id"
	indexPropertyBKCloudID         = "meta_bk_cloud_id"
	indexPropertyKeywords          = "keywords"
)

// elastic indexes.
var (
	indexBiz            *ESIndex
	indexSet            *ESIndex
	indexModule         *ESIndex
	indexHost           *ESIndex
	indexModel          *ESIndex
	indexObjectInstance *ESIndex
	indexList           []*ESIndex
)

func init() {
	// business application index.
	indexBiz = NewESIndex(indexNameBiz, indexVersionBiz, &ESIndexMetadata{
		Settings: ESIndexMetadataSettings{Shards: "1", Replicas: "1"},
		Mappings: ESIndexMetadataMappings{Properties: map[string]ESIndexMetadataMappingsProperty{
			indexPropertyID:                ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKObjID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKSupplierAccount: ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKBizID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyKeywords:          ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeText},
		}},
	})
	indexList = append(indexList, indexBiz)

	// set index.
	indexSet = NewESIndex(indexNameSet, indexVersionSet, &ESIndexMetadata{
		Settings: ESIndexMetadataSettings{Shards: "1", Replicas: "1"},
		Mappings: ESIndexMetadataMappings{Properties: map[string]ESIndexMetadataMappingsProperty{
			indexPropertyID:                ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKObjID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKSupplierAccount: ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKBizID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKParentID:        ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyKeywords:          ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeText},
		}},
	})
	indexList = append(indexList, indexSet)

	// module index.
	indexModule = NewESIndex(indexNameModule, indexVersionModule, &ESIndexMetadata{
		Settings: ESIndexMetadataSettings{Shards: "1", Replicas: "1"},
		Mappings: ESIndexMetadataMappings{Properties: map[string]ESIndexMetadataMappingsProperty{
			indexPropertyID:                ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKObjID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKSupplierAccount: ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKBizID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyKeywords:          ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeText},
		}},
	})
	indexList = append(indexList, indexModule)

	// host index.
	indexHost = NewESIndex(indexNameHost, indexVersionHost, &ESIndexMetadata{
		Settings: ESIndexMetadataSettings{Shards: "1", Replicas: "1"},
		Mappings: ESIndexMetadataMappings{Properties: map[string]ESIndexMetadataMappingsProperty{
			indexPropertyID:                ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKObjID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKSupplierAccount: ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKCloudID:         ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyKeywords:          ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeText},
		}},
	})
	indexList = append(indexList, indexHost)

	// model index.
	indexModel = NewESIndex(indexNameModel, indexVersionModel, &ESIndexMetadata{
		Settings: ESIndexMetadataSettings{Shards: "1", Replicas: "1"},
		Mappings: ESIndexMetadataMappings{Properties: map[string]ESIndexMetadataMappingsProperty{
			indexPropertyID:                ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKObjID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKSupplierAccount: ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKBizID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyKeywords:          ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeText},
		}},
	})
	indexList = append(indexList, indexModel)

	// object instance index.
	indexObjectInstance = NewESIndex(indexNameObjectInstance, indexVersionObjectInstance, &ESIndexMetadata{
		Settings: ESIndexMetadataSettings{Shards: "1", Replicas: "1"},
		Mappings: ESIndexMetadataMappings{Properties: map[string]ESIndexMetadataMappingsProperty{
			indexPropertyID:                ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKObjID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKSupplierAccount: ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyBKBizID:           ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeKeyword},
			indexPropertyKeywords:          ESIndexMetadataMappingsProperty{PropertyType: indexPropertyTypeText},
		}},
	})
	indexList = append(indexList, indexObjectInstance)

	log.Println("bk-cmdb elastic monstache plugin initialize successfully")
}

// ESIndexMetadataSettings elasticsearch index settings.
type ESIndexMetadataSettings struct {
	// Shards number of index shards as string type.
	Shards string `json:"number_of_shards"`

	// Replicas number of index document replicas as string type.
	Replicas string `json:"number_of_replicas"`
}

// ESIndexMetadataMappings elasticsearch index mappings.
type ESIndexMetadataMappings struct {
	// Properties elastic index properties.
	Properties map[string]ESIndexMetadataMappingsProperty `json:"properties"`
}

// ESIndexMetadataMappingsProperty elasticsearch index mappings property.
type ESIndexMetadataMappingsProperty struct {
	// PropertyType elastic index property type. Support 'keyword' 'text'.
	PropertyType string `json:"type"`
}

// ESIndexMetadata is elasticsearch index settings.
type ESIndexMetadata struct {
	// Settings elastic index settings.
	Settings ESIndexMetadataSettings `json:"settings"`

	// Mappings elastic index mappings.
	Mappings ESIndexMetadataMappings `json:"mappings"`
}

// ESIndex elasticsearch index.
type ESIndex struct {
	// name index name.
	name string

	// version is the plugin index version, as a postfix in target index.
	// the plugin would check and create the version index if it not exist,
	// and alias to bk-cmdb default index name.
	// NOTE: CHANGE the version name if you have modify the indexes metadata struct.
	version string

	// metadata index metadata including settings and mappings.
	metadata *ESIndexMetadata
}

// NewESIndex creates a new elasticsearch index.
func NewESIndex(name, version string, metadata *ESIndexMetadata) *ESIndex {
	return &ESIndex{name: name, version: version, metadata: metadata}
}

// Name returns the real elastic index name.
func (idx *ESIndex) Name() string {
	// return bk-cmdb.{index}_{version} as the real index name.
	return fmt.Sprintf("%s_%s", idx.name, idx.version)
}

// AliasName returns real bk-cmdb index name as index alias name.
func (idx *ESIndex) AliasName() string {
	return idx.name
}

// Metadata returns index metadata.
func (idx *ESIndex) Metadata() string {
	meta, err := ccjson.MarshalToString(idx.metadata)
	if err != nil {
		return ""
	}
	return meta
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

// analysisDocument analysis the given document, return document id and keywords.
func analysisDocument(document map[string]interface{}) (string, []string, error) {
	// analysis document id.
	documentID, ok := document["_id"].(string)
	if !ok {
		return "", nil, errors.New("missing document metadata id")
	}

	// analysis keywords.
	jsonDoc, err := ccjson.MarshalToString(document)
	if err != nil {
		return "", nil, err
	}

	// TODO: analysis in go types, not in json mode.
	keywords := analysisJSONKeywords(gjson.Parse(jsonDoc))

	// return document id and compressed keywords.
	return documentID, compressKeywords(keywords), nil
}

// indexingApplication indexing the business application instance.
func indexingApplication(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	documentID, keywords, err := analysisDocument(input.Document)
	if err != nil {
		return errors.New("missing document metadata id")
	}

	// build elastic document.
	document := map[string]interface{}{
		indexPropertyID:                documentID,
		indexPropertyBKObjID:           input.Document[common.BKObjIDField],
		indexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		indexPropertyBKBizID:           input.Document[common.BKAppIDField],
		indexPropertyKeywords:          keywords,
	}

	output.ID = documentID
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexBiz.AliasName()

	return nil
}

// indexingSet indexing the set instance.
func indexingSet(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	documentID, keywords, err := analysisDocument(input.Document)
	if err != nil {
		return errors.New("missing document metadata id")
	}

	// build elastic document.
	document := map[string]interface{}{
		indexPropertyID:                documentID,
		indexPropertyBKObjID:           input.Document[common.BKObjIDField],
		indexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		indexPropertyBKBizID:           input.Document[common.BKAppIDField],
		indexPropertyBKParentID:        input.Document[common.BKParentIDField],
		indexPropertyKeywords:          keywords,
	}

	output.ID = documentID
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexSet.AliasName()

	return nil
}

// indexingModule indexing the module instance.
func indexingModule(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	documentID, keywords, err := analysisDocument(input.Document)
	if err != nil {
		return errors.New("missing document metadata id")
	}

	// build elastic document.
	document := map[string]interface{}{
		indexPropertyID:                documentID,
		indexPropertyBKObjID:           input.Document[common.BKObjIDField],
		indexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		indexPropertyBKBizID:           input.Document[common.BKAppIDField],
		indexPropertyKeywords:          keywords,
	}

	output.ID = documentID
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexModule.AliasName()

	return nil
}

// indexingHost indexing the host instance.
func indexingHost(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	documentID, keywords, err := analysisDocument(input.Document)
	if err != nil {
		return errors.New("missing document metadata id")
	}

	// build elastic document.
	document := map[string]interface{}{
		indexPropertyID:                documentID,
		indexPropertyBKObjID:           input.Document[common.BKObjIDField],
		indexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		indexPropertyBKCloudID:         input.Document[common.BKCloudIDField],
		indexPropertyKeywords:          keywords,
	}

	output.ID = documentID
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
		return errors.New("missing model object id")
	}

	// query model.
	model := make(map[string]interface{})

	if err := input.MongoClient.Database(database).Collection(common.BKTableNameObjDes).
		FindOne(context.Background(), bson.D{{common.BKObjIDField, objectID}}).
		Decode(&model); err != nil {
		return fmt.Errorf("query model object[%s] failed, %+v", objectID, err)
	}

	// analysis document.
	documentID, keywords, err := analysisDocument(model)
	if err != nil {
		return errors.New("missing document metadata id")
	}

	// query model attribute.
	modelAttrs := []metadata.Attribute{}

	modelAttrsCursor, err := input.MongoClient.Database(database).Collection(common.BKTableNameObjAttDes).
		Find(context.Background(), bson.D{{common.BKObjIDField, objectID}})
	if err != nil {
		return fmt.Errorf("query model attributes object[%s] failed, %+v", objectID, err)
	}

	if err := modelAttrsCursor.All(context.Background(), &modelAttrs); err != nil {
		return fmt.Errorf("query model attributes object[%s] failed, %+v", objectID, err)
	}

	for _, attribute := range modelAttrs {
		jsonDoc, err := ccjson.MarshalToString(attribute)
		if err != nil {
			return fmt.Errorf("marshal model attributes object[%s] failed, %+v", objectID, err)
		}
		keywords = append(keywords, analysisJSONKeywords(gjson.Parse(jsonDoc))...)
	}

	// build elastic document.
	document := map[string]interface{}{
		indexPropertyID:                documentID,
		indexPropertyBKObjID:           objectID,
		indexPropertyBKSupplierAccount: model[common.BKOwnerIDField],
		indexPropertyBKBizID:           model[common.BKAppIDField],
		indexPropertyKeywords:          compressKeywords(keywords),
	}

	output.ID = documentID
	output.Document = document

	// use alias name to indexing document.
	output.Index = indexModel.AliasName()

	return nil
}

// indexingObjectInstance indexing the common object instance.
func indexingObjectInstance(input *monstachemap.MapperPluginInput, output *monstachemap.MapperPluginOutput) error {
	// analysis document.
	documentID, keywords, err := analysisDocument(input.Document)
	if err != nil {
		return errors.New("missing document metadata id")
	}

	// build elastic document.
	document := map[string]interface{}{
		indexPropertyID:                documentID,
		indexPropertyBKObjID:           input.Document[common.BKObjIDField],
		indexPropertyBKSupplierAccount: input.Document[common.BKOwnerIDField],
		indexPropertyBKBizID:           input.Document[common.BKAppIDField],
		indexPropertyKeywords:          keywords,
	}

	output.ID = documentID
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
