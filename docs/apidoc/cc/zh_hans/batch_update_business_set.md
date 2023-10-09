### 功能描述

更新业务集信息(v3.10.12+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_ids | array  | 是  | 业务集ID列表 |
| data           | object | 是 | 业务集数据 |

#### data

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_attr |  object  | 否     | 业务集模型字段 |
| bk_scope  |  object  | 否     | 所选业务范围 |

#### bk_biz_set_attr

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_name   |  string  | 是     | 业务集名称|
| bk_biz_maintainer |  string  | 否     | 运维人员 |
| bk_biz_set_desc   |  string  | 否     | 业务集描述 |

#### bk_scope

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| match_all |  bool  | 是     | 所选业务范围标记|
| filter    |  object| 否     | 所选业务的范围条件 |

#### filter

该参数为业务属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合仅支持AND操作，可以嵌套，最多嵌套2层。

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
- 此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段
- 对于批量场景(bk_biz_set_ids中的ID数量大于1)下不允许更改`bk_biz_set_name`和`bk_scope`字段
- 每次批量更新最大数量为200。

### 请求参数示例

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_ids":[
        2
    ],
    "data":{
        "bk_biz_set_attr":{
            "bk_biz_set_name": "test",
            "bk_biz_set_desc":"xxx",
            "biz_set_maintainer":"xxx"
        },
        "bk_scope":{
            "match_all":false,
            "filter":{
                "condition":"AND",
                "rules":[
                    {
                        "field":"bk_sla",
                        "operator":"equal",
                        "value":"2"
                    }
                ]
            }
        }
    }
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission":null,
    "data": {},
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
