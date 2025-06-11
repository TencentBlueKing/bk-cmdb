### 描述

查询业务集(版本：v3.10.12+，权限：业务集查看权限)

### 输入参数

| 参数名称              | 参数类型   | 必选 | 描述                           |
|-------------------|--------|----|------------------------------|
| bk_biz_set_filter | object | 否  | 业务集条件范围                      |
| time_condition    | object | 否  | 业务集时间范围                      |
| fields            | array  | 否  | 查询条件，参数为业务的任意属性，如果不写代表搜索全部数据 |
| page              | object | 是  | 分页条件                         |

#### bk_biz_set_filter

该参数为业务集属性字段过滤规则的组合，用于根据业务集属性字段搜索业务集。组合支持AND 和 OR 两种方式，允许嵌套，最多嵌套2层。

| 参数名称      | 参数类型   | 必选 | 描述        |
|-----------|--------|----|-----------|
| condition | string | 是  | 规则操作符     |
| rules     | array  | 是  | 过滤业务的范围规则 |

#### rules

过滤规则为三元组 `field`, `operator`, `value`

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### time_condition

| 参数名称  | 参数类型   | 必选 | 描述           |
|-------|--------|----|--------------|
| oper  | string | 是  | 操作符，目前只支持and |
| rules | array  | 是  | 时间查询条件       |

#### rules

| 参数名称  | 参数类型   | 必选 | 描述                          |
|-------|--------|----|-----------------------------|
| field | string | 是  | 取值为模型的字段名                   |
| start | string | 是  | 起始时间，格式为yyyy-MM-dd hh:mm:ss |
| end   | string | 是  | 结束时间，格式为yyyy-MM-dd hh:mm:ss |

#### page

| 参数名称         | 参数类型   | 必选 | 描述                                                        |
|--------------|--------|----|-----------------------------------------------------------|
| start        | int    | 是  | 记录开始位置                                                    |
| limit        | int    | 是  | 每页限制条数,最大500                                              |
| enable_count | bool   | 是  | 本次请求是否为获取数量还是详情的标记                                        |
| sort         | string | 否  | 排序字段，通过在字段前面增加 -，如 sort:&#34;-field&#34; 可以表示按照字段 field降序 |

**注意：**

- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- `sort`如果调用方没有指定，后台默认指定为业务集ID。

### 调用示例

```json
{
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

### 响应示例

#### 详细信息接口响应

```json
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
}
```

#### 业务集数量接口响应

```json
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
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称  | 参数类型  | 描述     |
|-------|-------|--------|
| count | int   | 记录条数   |
| info  | array | 业务实际数据 |

#### info

| 参数名称              | 参数类型   | 描述        |
|-------------------|--------|-----------|
| bk_biz_set_id     | int    | 业务集ID     |
| create_time       | string | 业务集创建时间   |
| last_time         | string | 业务集修改时间   |
| bk_biz_set_name   | string | 业务集名称     |
| bk_biz_maintainer | string | 运维人员      |
| bk_biz_set_desc   | string | 业务集描述     |
| bk_scope          | object | 业务集所选业务范围 |
| bk_created_at     | string | 创建时间      |
| bk_created_by     | string | 创建人       |
| bk_updated_at     | string | 更新时间      |

#### bk_scope

| 参数名称      | 参数类型   | 描述        |
|-----------|--------|-----------|
| match_all | bool   | 所选业务范围标记  |
| filter    | object | 所选业务的范围条件 |

#### filter

该参数为业务属性字段过滤规则的组合，用于根据业务属性字段搜索业务。组合仅支持AND操作，可以嵌套，最多嵌套2层。

| 参数名称      | 参数类型   | 描述          |
|-----------|--------|-------------|
| condition | string | 规则操作符       |
| rules     | array  | 所选业务的范围条件规则 |

#### rules

| 参数名称     | 参数类型   | 描述                          |
|----------|--------|-----------------------------|
| field    | string | 字段名                         |
| operator | string | 操作符,可选值 equal,in            |
| value    | -      | 操作数,不同的operator对应不同的value格式 |

**注意：**

- 如果本次请求是查询详细信息那么count为0，如果查询的是数量，那么info为空。
- 此处的输入针对`info`参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段
