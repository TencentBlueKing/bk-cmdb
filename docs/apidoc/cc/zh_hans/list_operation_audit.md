### 功能描述

 根据条件获取操作审计日志

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型      | 必选   |  描述                       |
|---------------------|------------|--------|-----------------------------|
| page                | object     | 是     | 分页参数                    |
| condition           | object     | 否     | 操作审计日志查询条件                   |

#### page

| 字段      |  类型      | 必选   |  描述                |
|-----------|------------|--------|----------------------|
| start     |  int       | 否     | 记录开始位置         |
| limit     |  int       | 是     | 每页限制条数,最大200 |
| sort      |  string    | 否     | 排序字段             |

#### condition

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id     |int      |否      | 业务id                                               |
| resource_type  |string      |否      | 操作的具体资源类型 |
| action     |    array  |  否    |  操作类型 |
|   operation_time   |    object  |  是    | 操作时间 |
|   user   |    string  |  否    | 操作人 |
|    resource_name  |    string  |  否    | 资源名称 |
|    category  |    string  |  否    | 查询的类型 |
| fuzzy_query    | bool         | 否   | 是否使用模糊查询对资源名称进行查询，**模糊查询效率低，性能差**，该字段仅对resource_name产生影响，使用condition方式进行模糊查询时会忽略该字段，请二者选其一使用。 |
| condition | array | 否 | 指定查询条件，与user和resource_name不能同时提供 |

##### condition.condition

| 字段     | 类型         | 必选 | 描述                                                         |
| -------- | ------------ | ---- | ------------------------------------------------------------ |
| field    | string       | 是   | 对象的字段，仅为"user"，"resource_name"                      |
| operator | string       | 是   | 操作符，in 为属于，not_in 为不属于, contains 为包含,field为resource_name时可以使用contains进行模糊查询 |
| value    | string/array | 是   | 字段对应的值，in和not_in需要array类型，contains需要string类型 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition":{
        "bk_biz_id":2,
        "resource_type":"host",
        "action":[
            "create",
            "delete"
        ],
        "operation_time":{
            "start":"2020-09-23 00:00:00",
            "end":"2020-11-01 23:59:59"
        },
        "user":"admin",
        "resource_name":"1.1.1.1",
        "category":"host",
        "fuzzy_query": false
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"-operation_time"
    }
}
```

```json
{
    "condition":{
        "bk_biz_id":2,
        "resource_type":"host",
        "action":[
            "create",
            "delete"
        ],
        "operation_time":{
            "start":"2020-09-23 00:00:00",
            "end":"2020-11-01 23:59:59"
        },
      	"condition":[
          {
            "field":"user",
            "operator":"in",
            "value":["admin"]
          },
          {
            "field":"resource_name",
            "operator":"in",
            "value":["1.1.1.1"]
          }
        ],
        "category":"host"
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"-operation_time"
    }
}
```

### 返回结果示例

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data":{
        "count":2,
        "info":[
            {
                "id":7,
                "audit_type":"",
                "bk_supplier_account":"",
                "user":"admin",
                "resource_type":"host",
                "action":"delete",
                "operate_from":"",
                "operation_detail":null,
                "operation_time":"2020-10-09 21:30:51",
                "bk_biz_id":1,
                "resource_id":4,
                "resource_name":"2.2.2.2"
            },
            {
                "id":2,
                "audit_type":"",
                "bk_supplier_account":"",
                "user":"admin",
                "resource_type":"host",
                "action":"delete",
                "operate_from":"",
                "operation_detail":null,
                "operation_time":"2020-10-09 17:13:55",
                "bk_biz_id":1,
                "resource_id":1,
                "resource_name":"1.1.1.1"
            }
        ]
    }
}
```

### 返回结果参数说明
#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | object | 请求返回的数据 |

#### data

| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
| count     | int       | 记录条数     |
| info      | array     | 操作审计的记录信息 |

#### info
| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
|    id |      int  |    审计ID  |
|   audit_type  |     string   |   操作审计类型   |
|   bk_supplier_account  |    string    | 开发商账号     |
|   user  |      string  |    操作人  |
|   resource_type  |    string    |   资源类型   |
|  action   |    string    |    操作类型  |
|    operate_from |    string    |   来源平台   |
|  operation_detail   |     object     |  操作细节    |
| operation_time    |     string   |    操作时间  |
|  bk_biz_id   |       int |    业务id  |
| resource_id    |     int   |    资源id  |
|   resource_name  |     string   |  资源名称    |