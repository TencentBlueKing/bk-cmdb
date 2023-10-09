### 功能描述

查询业务集中的业务(v3.10.12+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_id | int    | 是     | 业务集ID |
| filter      |  object  | 否     | 业务属性组合查询条件 |
| fields      |  array   | 否     | 指定查询的字段，参数为业务的任意属性，如果不填写字段信息，系统会返回业务的所有字段 |
| page        |  object  | 是     | 分页条件 |

#### filter

查询条件。组合支持AND 和 OR 两种方式。可以嵌套，最多嵌套2层。

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

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start    |  int    | 是     | 记录开始位置 |
| limit    |  int    | 是     | 每页限制条数,最大500 |
| enable_count |  bool  | 是  | 是否获取查询对象数量的标记 |
| sort     |  string | 否     | 排序字段，通过在字段前面增加 -，如 sort:&#34;-field&#34; 可以表示按照字段 field降序 |

**注意：**
- `enable_count` 如果此标记为true那么表示此次请求是获取数量。此时其余字段必须为初始化值,start为0,limit为:0, sort为""。

### 请求参数示例

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_id":2,
    "filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"xxx",
                "operator":"equal",
                "value":"xxx"
            },
            {
                "field":"xxx",
                "operator":"in",
                "value":[
                    "xxx"
                ]
            }
        ]
    },
    "fields":[
        "bk_biz_id",
        "bk_biz_name"
    ],
    "page":{
        "start":0,
        "limit":10,
        "enable_count":false,
        "sort":"bk_biz_id"
    }
}
```

### 详细信息返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission":null,
    "data": {
        "count": 0,
        "info": [
            {
                "bk_biz_id": 1,
                "bk_biz_name": "xxx"
            }
        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### 查询业务数量的返回结果示例

```python
{
    "result":true,
    "code":0,
    "message":"",
    "permission":null,
    "data":{
        "count":10,
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

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| bk_biz_id     | int       | 业务ID |
| bk_biz_name      | string     | 业务名称 |


**注意：**
- 如果本次请求是查询详细信息那么count为0，如果查询的是数量，那么info为空。
