### 功能描述

根据关联关系实例查询模型实例

- 该接口只适用于自定义层级模型和通用模型实例上，不适用于业务、集群、模块、主机等模型实例

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型      | 必选   |  描述                       |
|---------------------|------------|--------|-----------------------------|
| bk_obj_id           | string     | 是     | 模型ID                      |
| page                | object     | 是     | 分页参数                    |
| condition           | object     | 否     | 具有关联关系的模型实例查询条件                    |
| time_condition      | object     | 否     | 按时间查询模型实例的查询条件 |
| fields              | object     | 否     | 指定查询模型实例返回的字段,key为模型ID，value为该查询模型要返回的模型属性字段|

#### page

| 字段      |  类型      | 必选   |  描述                |
|-----------|------------|--------|----------------------|
| start     |  int       | 是     | 记录开始位置         |
| limit     |  int       | 是     | 每页限制条数,最大200 |
| sort      |  string    | 否     | 排序字段             |

#### condition
示例中的user为模型

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| field     |string      |是      | 取值为模型的字段名                                               |
| operator  |string      |是      | 取值为：$regex $eq $ne                                           |
| value     |string      |是      | field配置的模型字段名所对应的值                                  |          

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


### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "bk_switch",
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_inst_id"
    },
    "fields": {
        "bk_switch": [
            "bk_asset_id",
            "bk_inst_id",
            "bk_inst_name",
            "bk_obj_id"
        ]
    },
    "condition": {
        "user": [
            {
                "field": "operator",
                "operator": "$regex",
                "value": "admin"
            }
        ]
    },
    "time_condition": {
        "oper": "and",
        "rules": [
            {
                "field": "create_time",
                "start": "2021-05-13 01:00:00",
                "end": "2021-05-14 01:00:00"
            }
        ]
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 2,
        "info": [
            {
                "bk_asset_id": "sw00001",
                "bk_inst_id": 1,
                "bk_inst_name": "sw1",
                "bk_obj_id": "bk_switch"
            },
            {
                "bk_asset_id": "sw00002",
                "bk_inst_id": 2,
                "bk_inst_name": "sw2",
                "bk_obj_id": "bk_switch"
            }
        ]
    }
}
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

#### data

| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
| count     | int       | 记录条数     |
| info      | array     | 模型实例实际数据 |
