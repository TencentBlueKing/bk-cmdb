### 功能描述

查询业务集(v3.10.12+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_filter | object  | 否   | 业务集条件范围 |
| time_condition    | object  | 否   | 业务集时间范围 |
| fields            | array   | 否   | 查询条件，参数为业务的任意属性，如果不写代表搜索全部数据 |
| page              | object  | 是   | 分页条件 |

#### bk_biz_set_filter

该参数为业务集属性字段过滤规则的组合，用于根据业务集属性字段搜索业务集。组合支持AND 和 OR 两种方式，允许嵌套，最多嵌套2层。

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition |  string  | 是    | 规则操作符|
| rules |  array  | 是     | 过滤业务的范围规则 |


#### rules
过滤规则为三元组 `field`, `operator`, `value`

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| -------- | ------ | ---- | ------ | ------ | ------------------------------------------------------------ |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### time_condition

| 字段   | 类型   | 必选 |  描述              |
|-------|--------|-----|--------------------|
| oper  | string | 是  | 操作符，目前只支持and |
| rules | array  | 是  | 时间查询条件         |

#### rules

| 字段   | 类型   | 必选 | 描述                             |
|-------|--------|-----|----------------------------------|
| field | string | 是  | 取值为模型的字段名                  |
| start | string | 是  | 起始时间，格式为yyyy-MM-dd hh:mm:ss |
| end   | string | 是  | 结束时间，格式为yyyy-MM-dd hh:mm:ss | 

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大500 |
| enable_count |  bool | 是 | 本次请求是否为获取数量还是详情的标记 |
| sort     |  string | 否     | 排序字段，通过在字段前面增加 -，如 sort:&#34;-field&#34; 可以表示按照字段 field降序 |

**注意：**
- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- `sort`如果调用方没有指定，后台默认指定为业务集ID。

### 请求参数示例

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"bk_biz_set_id",
                "operator":"equal",
                "value":10
            },
            {
                "field":"bk_biz_maintainer",
                "operator":"equal",
                "value":"admin"
            }
        ]
    },
    "time_condition":{
        "oper":"and",
        "rules":[
            {
                "field":"create_time",
                "start":"2021-05-13 01:00:00",
                "end":"2021-05-14 01:00:00"
            }
        ]
    },
    "fields": [
        "bk_biz_id"
    ],
    "page":{
        "start":0,
        "limit":500,
        "enable_count":false,
        "sort":"bk_biz_set_id"
    }
}
```

### 返回结果示例

### 详细信息接口响应
```python

{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":0,
        "info":[
            {
                "bk_biz_set_id":10,
                "bk_biz_set_name":"biz_set",
                "bk_biz_set_desc":"dba",
                "bk_biz_maintainer":"tom",
                "create_time":"2021-09-06T08:10:50.168Z",
                "last_time":"2021-10-15T02:30:01.867Z",
                "bk_scope":{
                    "match_all":true
                }
            },
            {
                "bk_biz_set_id":11,
                "bk_biz_set_name":"biz_set1",
                "bk_biz_set_desc":"dba",
                "bk_biz_maintainer":"tom",
                "create_time":"2021-09-06T08:10:50.168Z",
                "last_time":"2021-10-15T02:30:01.867Z",
                "bk_scope":{
                    "match_all":false,
                    "filter":{
                        "condition":"AND",
                        "rules":[
                            {
                                "field":"bk_sla",
                                "operator":"equal",
                                "value":"3"
                            },
                            {
                                "field":"bk_biz_maintainer",
                                "operator":"equal",
                                "value":"admin"
                            }
                        ]
                    }
                }
            }
        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### 业务集数量接口响应
```python
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":2,
        "info":[
        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| data    | object | 请求返回的数据                           |
| request_id    | string | 请求链id    |

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int       | 记录条数 |
| info      | array     | 业务实际数据 |

#### info

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_id   |  int  | 是   | 业务集ID|
| create_time   |  string  | 否   | 业务集创建时间|
| last_time   |  string  | 否   | 业务集修改时间|
| bk_biz_set_name   |  string  | 是   | 业务集名称|
| bk_biz_maintainer |  string  | 否   | 运维人员 |
| bk_biz_set_desc   |  string  | 否   | 业务集描述 |
| bk_scope   |  object  | 否   | 业务集所选业务范围 |

#### bk_scope

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| match_all |  bool  | 是    | 所选业务范围标记|
| filter |  object  | 否     | 所选业务的范围条件 |

#### filter

该参数为业务属性字段过滤规则的组合，用于根据业务属性字段搜索业务。组合仅支持AND操作，可以嵌套，最多嵌套2层。

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| condition |  string  | 是    | 规则操作符|
| rules |  array  | 是     | 所选业务的范围条件规则 |


#### rules

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| -------- | ------ | ---- | ------ | ------ | ------------------------------------------------------------ |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,in |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

**注意：**
- 如果本次请求是查询详细信息那么count为0，如果查询的是数量，那么info为空。
- 此处的输入针对`info`参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段
