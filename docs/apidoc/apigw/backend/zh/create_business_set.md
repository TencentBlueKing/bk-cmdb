### 描述

新建业务集(版本：v3.10.12+，权限：业务集新增权限)

### 输入参数

| 参数名称            | 参数类型   | 必选 | 描述      |
|-----------------|--------|----|---------|
| bk_biz_set_attr | object | 是  | 业务集模型字段 |
| bk_scope        | object | 是  | 所选业务范围  |

#### bk_biz_set_attr

| 参数名称              | 参数类型   | 必选 | 描述    |
|-------------------|--------|----|-------|
| bk_biz_set_name   | string | 是  | 业务集名称 |
| bk_biz_maintainer | string | 否  | 运维人员  |
| bk_biz_set_desc   | string | 否  | 业务集描述 |

#### bk_scope

| 参数名称      | 参数类型   | 必选 | 描述        |
|-----------|--------|----|-----------|
| match_all | bool   | 是  | 所选业务范围标记  |
| filter    | object | 否  | 所选业务的范围条件 |

#### filter

该参数为业务属性字段过滤规则的组合，用于根据业务属性字段搜索业务。组合仅支持AND操作，可以嵌套，最多嵌套2层。

| 参数名称      | 参数类型   | 必选 | 描述          |
|-----------|--------|----|-------------|
| condition | string | 是  | 规则操作符       |
| rules     | array  | 是  | 所选业务的范围条件规则 |

#### rules

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

**注意：**

- 此处的输入针对`bk_biz_set_attr`参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段
- `bk_scope`中`match_all`字段是true的时候表示业务集的所选业务范围是全部，此时不需要填写参数`filter`。如果`match_all`
  字段是false，`filter`需要非空，用户需要显式的指
  定业务的选择范围
- 业务集中所选业务属性圈定类型是组织和枚举

### 调用示例

```json
{
    "bk_biz_set_attr":{
        "bk_biz_set_name":"biz_set",
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
                    "value":"3"
                },
                {
                    "field":"life_cycle",
                    "operator":"equal",
                    "value":1
                }
            ]
        }
    }
}
```

### 响应示例

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":5,
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | int    | 创建的业务集id                   |
