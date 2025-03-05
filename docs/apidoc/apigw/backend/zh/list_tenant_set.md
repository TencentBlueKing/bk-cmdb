### 描述

查询租户集(版本：v3.15.1+，权限：租户集查看权限)

### 输入参数

| 参数名称   | 参数类型   | 必选 | 描述                                |
|--------|--------|----|-----------------------------------|
| filter | object | 否  | 租户集查询条件                           |
| fields | array  | 否  | 租户集属性列表，控制返回结果里有哪些字段，如果不写代表搜索全部数据 |
| page   | object | 是  | 分页条件                              |

#### filter 字段说明

租户集的属性字段过滤规则，用于根据租户集的属性字段搜索数据。该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

##### 组合过滤规则

由其它规则组合而成的过滤规则，组合的规则间支持逻辑与/或关系

| 参数名称      | 参数类型   | 必选 | 描述                              |
|-----------|--------|----|---------------------------------|
| condition | string | 是  | 组合查询条件，支持 `AND` 和 `OR` 两种方式     |
| rules     | array  | 是  | 查询规则，可以是 `组合查询参数` 或 `原子查询参数` 类型 |

##### 原子过滤规则

基础的过滤规则，表示对某一个字段进行过滤的规则。任何过滤规则都直接是原子过滤规则, 或由多个原子过滤规则组合而成

| 参数名称     | 参数类型                          | 必选 | 描述                                                                                                |
|----------|-------------------------------|----|---------------------------------------------------------------------------------------------------|
| field    | string                        | 是  | 租户集的字段                                                                                            |
| operator | string                        | 是  | 操作符，可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | 不同的field和operator对应不同的value格式 | 否  | 操作值                                                                                               |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md>

#### page 字段说明

| 参数名称         | 参数类型   | 必选 | 描述                                                                           |
|--------------|--------|----|------------------------------------------------------------------------------|
| start        | int    | 是  | 记录开始位置                                                                       |
| limit        | int    | 是  | 每页限制条数，最大500                                                                 |
| sort         | string | 否  | 排序字段                                                                         |
| enable_count | bool   | 是  | 是否获取查询对象数量的标记。如果此标记为true那么表示此次请求是获取数量，此时其余字段必须为初始化值，start为0，limit为:0，sort为"" |

**注意：**

- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- 必须设置分页参数，一次最大查询数据不超过500个。

### 调用示例

```json
{
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "id",
        "operator": "equal",
        "value": 1
      }
    ]
  },
  "fields": [
    "name"
  ],
  "page": {
    "start": 0,
    "limit": 500,
    "enable_count": false
  }
}
```

### 响应示例

#### 详细信息接口响应

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 0,
    "info": [
      {
        "id": 1,
        "name": "All tenants",
        "maintainer": "",
        "description": "全租户",
        "default": 1,
        "bk_scope": {
          "match_all": true
        },
        "bk_created_at": "2025-03-03 10:00:00",
        "bk_created_by": "admin",
        "bk_updated_at": "2025-03-03 10:00:00",
        "bk_updated_by": "admin"
      }
    ]
  }
}
```

#### 租户集数量接口响应

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 2,
    "info": [
    ]
  }
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

| 参数名称  | 参数类型  | 描述    |
|-------|-------|-------|
| count | int   | 记录条数  |
| info  | array | 租户集数据 |

#### info

| 参数名称          | 参数类型   | 描述                |
|---------------|--------|-------------------|
| id            | int    | 租户集ID             |
| name          | string | 租户集名称             |
| maintainer    | string | 运维人员              |
| description   | string | 租户集描述             |
| default       | int    | 租户集类型，1表示内置全租户租户集 |
| bk_scope      | object | 租户集中所选租户范围        |
| bk_created_at | string | 创建时间              |
| bk_created_by | string | 创建人               |
| bk_updated_at | string | 更新时间              |
| bk_updated_by | string | 更新人               |

#### bk_scope

| 参数名称      | 参数类型 | 描述       |
|-----------|------|----------|
| match_all | bool | 是否匹配全部租户 |

