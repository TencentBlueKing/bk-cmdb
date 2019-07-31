# 全文检索

## 方案
基于强大的elasticsearch进行搜索，把cmdb的mongo数据使用mongo-connector同步到
elasticsearch，封装es的全文搜索的API提供出来。
mongo-connector，通过读取mongodb的replica oplog，将mongodb产生的操作
在elasticsearch上replay，来实现单向同步。即mongo里有数据变动，mongo-connector
就会把相应的数据同步到es中同时进行更新。

[原理图](../resource/img/mongo-connector.png)

## es的使用
全文检索使用了es的query_string的参数，并且配合使用bool(must, must_not, should)
进行了定制化的搜索，配合使用aggs进行了数据的分类汇聚上报，配合使用了highlight提供了
数据的高亮。
一个完整的es的query请求大致如下：
\*e\*为搜索条件，
```
{
    "query": {
        "bool": {
            "must": [
                {
                    "term": { "bk_obj_id": "test_search"}
                },
                {
                    "query_string": {"query": "*e*"}
                }
            ],
            "must_not": [
                {
                    "match": {"bk_supplier_account": "*e*"}
                }
            ],
            "should": [
                {
                    "bool": {
                        "must_not": [
                            {
                                "regexp": { "metadata.label.bk_biz_id": "[0-9]*" }
                            }
                        ]
                    }
                },
                {
                    "term": { "metadata.label.bk_biz_id": "2" }
                }
            ],
            "minimum_should_match" : 1
        }
    },
    "aggs": {
        "bk_obj_id_agg": {
            "terms": {
                "field": "bk_obj_id.keyword"
            }
        },
        "type_agg": {
            "terms": {
                "field": "_type"
            }
        }
    },
    "highlight": {
        "fields": {
            "*" : {}
        },
        "require_field_match":  false
    }
}
```
把返回的结果封装转换成cmdb的api返回值规范：
[全文检索api](../apidoc/v3.5/full_text_find.md)

## mongo-connector和es的部署
[部署](../overview/installation.md)
第6和第7步，以及后面的配置开关full_text_search(值为off或者on)

## 参考github
[olivere elastic](https://github.com/olivere/elastic)

[mongo-connector](https://github.com/yougov/mongo-connector)


## 参考wiki
[olivere elastic wiki](https://github.com/yougov/mongo-connector/wiki/Usage-with-ElasticSearch)

[mongo-connector wiki](https://github.com/olivere/elastic/wiki)