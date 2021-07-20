
### 全文检索
* API:  POST /api/v3/find/full_text
* API名称： full_text_find
* 功能说明：
	* 中文：全文检索
	* English ：full text find
* input body:
```
{
    "page": {
        "start": 0,
        "limit": 10
    },
    "query_string": "wtop",
    "bk_obj_id": "test_search",
    "bk_biz_id": "2",
    "filter": ["model"]
}
```


* input字段说明：

| 名称         | 类型        | 必填 | 默认值 | 说明       | Description                                                   |
| ------------ | ----------- | ---- | ------ | ---------- | ------------------------------------------------------------- |
| start        | int         | 否   | -1     | 分页开始   | 分页开始                                                      |
| limit        | int         | 否   | -1     | 分页大小   | 分页大小                                                      |
| query_string | string      | 否   | 无     | 检索内容   | query string content                                          |
| bk_obj_id    | string      | 否   | 无     | 对象实例id | 指定对象                                                      |
| bk_biz_id    | string      | 否   | 无     | 业务id     | 指定业务                                                      |
| filter       | string list | 否   | 空     | 过滤搜索表 | 取值范围["model", "object", "host", "process", "application"] |


* output:

```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": {
        "total": 1,
        "aggregations": [
            {
                "key": "test_search",
                "count": 1
            }
        ],
        "hits": [
            {
                "source": {
                    "bk_inst_id": 19,
                    "bk_inst_name": "liow",
                    "bk_obj_id": "test_search",
                    "bk_supplier_account": "0",
                    "bool_e": false,
                    "char_a": "123wdw",
                    "char_f": "awe",
                    "enum_d": "1",
                    "float_c": 23,
                    "int_b": 123,
                    "longchar_g": "wtop",
                    "metadata": {
                        "label": {
                            "bk_biz_id": "2"
                        }
                    }
                },
                "highlight": {
                    "bk_obj_id": [
                        "<em>test_search</em>"
                    ],
                    "bk_obj_id.keyword": [
                        "<em>test_search</em>"
                    ],
                    "longchar_g": [
                        "<em>wtop</em>"
                    ],
                    "longchar_g.keyword": [
                        "<em>wtop</em>"
                    ],
                    "metadata.label.bk_biz_id": [
                        "<em>2</em>"
                    ],
                    "metadata.label.bk_biz_id.keyword": [
                        "<em>2</em>"
                    ]
                },
                "type": "object",
                "score": 2.2249658
            }
        ]
    }
}
```
*  output字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

data说明：

| 名称         | 类型        | 说明         | Description                |
| ------------ | ----------- | ------------ | -------------------------- |
| total        | int         | 总数         | 搜索匹配的总数             |
| aggregations | object list | 数据汇聚     | 搜索的匹配结果分类统计     |
| key          | string      | 模型id       | 匹配结果的所属模型id       |
| count        | int         | 计数         | 匹配结果中所属此模型的数量 |
| hits         | object list | 匹配结果列表 | 搜索的匹配结果的集合       |
| type         | string      | 类型         | 查询结果所属数据类型       |
| score        | float       | 最佳匹配得分 | 搜索结果匹配程度           |
| source       | object      | 属性值       | 搜索结果的具体内容         |
| highlight    | object      | 高亮字段     | 匹配高亮显示的字段         |
