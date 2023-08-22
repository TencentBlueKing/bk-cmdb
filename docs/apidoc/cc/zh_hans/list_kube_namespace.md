### 功能描述

查询namespace (版本：v3.12.1+，权限：业务访问)

### 请求参数

{{ common_args_desc }}

#### 接口参数

- 通用字段：

| 字段        | 类型     | 必选  | 描述                                  |
|-----------|--------|-----|-------------------------------------|
| bk_biz_id | int    | 是   | 业务id                                |
| filter    | object | 否   | namespace查询条件                       |
| fields    | array  | 否   | 属性列表，控制返回结果里有哪些字段，能够加速接口请求和减少网络流量传输 |
| page      | object | 是   | 分页信息                                |

#### filter 字段说明

namespace的属性字段过滤规则，用于根据namespace的属性字段搜索数据。该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

##### 组合过滤规则

由其它规则组合而成的过滤规则，组合的规则间支持逻辑与/或关系

| 字段        | 类型     | 必选  | 描述                              |
|-----------|--------|-----|---------------------------------|
| condition | string | 是   | 组合查询条件，支持 `AND` 和 `OR` 两种方式     |
| rules     | array  | 是   | 查询规则，可以是 `组合过滤规则` 或 `原子过滤规则` 类型 |

##### 原子过滤规则

基础的过滤规则，表示对某一个字段进行过滤的规则。任何过滤规则都直接是原子过滤规则, 或由多个原子过滤规则组合而成

| 名称       | 类型                            | 必选  | 说明                                                                                                |
|----------|-------------------------------|-----|---------------------------------------------------------------------------------------------------|
| field    | string                        | 是   | container的字段                                                                                      |
| operator | string                        | 是   | 操作符，可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between | 
| value    | 不同的field和operator对应不同的value格式 | 否   | 操作值                                                                                               |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md>

#### page

| 字段           | 类型     | 必选  | 描述                                                                           |
|--------------|--------|-----|------------------------------------------------------------------------------|
| start        | int    | 是   | 记录开始位置                                                                       |
| limit        | int    | 是   | 每页限制条数，最大500                                                                 |
| sort         | string | 否   | 排序字段                                                                         |
| enable_count | bool   | 是   | 是否获取查询对象数量的标记。如果此标记为true那么表示此次请求是获取数量，此时其余字段必须为初始化值，start为0，limit为:0，sort为"" |

**注意：**

- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- 必须设置分页参数，一次最大查询数据不超过500个。

### 请求参数示例

### 获取详细信息请求参数

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "bk_cluster_id",
        "operator": "equal",
        "value": 1
      },
      {
        "field": "name",
        "operator": "equal",
        "value": "test"
      }
    ]
  },
  "fields": [
    "name"
  ],
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "name",
    "enable_count": false
  }
}
```

### 获取数量请求示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "bk_cluster_id",
        "operator": "equal",
        "value": 1
      },
      {
        "field": "name",
        "operator": "equal",
        "value": "test"
      }
    ]
  },
  "page": {
    "enable_count": true
  }
}
```

### 返回结果示例

### 详细信息接口响应

```json

{
  "result": true,
  "code": 0,
  "data": {
    "count": 0,
    "info": [
      {
        "name": "test"
      }
    ]
  },
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 获取数量返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 100,
    "info": [
    ]
  },
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 返回结果参数说明

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

#### data

| 字段    | 类型    | 描述                    |
|-------|-------|-----------------------|
| count | int   | 记录条数                  |
| info  | array | 实际数据，仅返回fields里设置了的字段 |

#### info[x]

| 字段              | 类型     | 描述                         |
|-----------------|--------|----------------------------|
| name            | string | 命名空间名称                     |
| labels          | map    | 标签                         |
| resource_quotas | array  | 命名空间CPU与内存的requests与limits |

#### resource_quotas[x]

| 字段             | 类型     | 描述                                                                                                                 |
|----------------|--------|--------------------------------------------------------------------------------------------------------------------|
| hard           | object | 每个命名资源所需的硬限制                                                                                                       |
| scopes         | array  | 配额作用域,可选值为："Terminating"、"NotTerminating"、"BestEffort"、"NotBestEffort"、"PriorityClass"、"CrossNamespacePodAffinity" |
| scope_selector | object | 作用域选择器                                                                                                             |

#### scope_selector

| 字段                | 类型    | 描述    |
|-------------------|-------|-------|
| match_expressions | array | 匹配表达式 |

#### match_expressions[x]

| 字段         | 类型     | 描述                                                                                                                 |
|------------|--------|--------------------------------------------------------------------------------------------------------------------|
| scope_name | array  | 配额作用域,可选值为："Terminating"、"NotTerminating"、"BestEffort"、"NotBestEffort"、"PriorityClass"、"CrossNamespacePodAffinity" |
| operator   | string | 选择器操作符，可选值为："In"、"NotIn"、"Exists"、"DoesNotExist"                                                                   |
| values     | array  | 字符串数组，如果操作符为"In"或"NotIn",不能为空，如果为"Exists"或"DoesNotExist"，必须为空                                                      |

**注意：**

- 如果本次请求是查询详细信息那么count为0，如果查询的是数量，那么info为空
