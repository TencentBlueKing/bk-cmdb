### 描述

查询集群

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述                                        |
|---------------------|--------|----|-------------------------------------------|
| bk_supplier_account | string | 否  | 开发商账号                                     |
| bk_biz_id           | int    | 是  | 业务id                                      |
| fields              | array  | 是  | 查询字段，所有字段均为set定义的字段，这些字段包括预置字段，也包括用户自定义字段 |
| condition           | dict   | 是  | 查询条件，所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段 |
| filter| object| 否  | 属性组合查询条件                                                         |
| time_condition | object | 否  | 按时间查询模型实例的查询条件  |
| page                | dict   | 是  | 分页条件                                      |

- `filter`与`condition`两个参数只能有一个生效，参数`condition`不建议继续使用。
- 参数`filter` 中涉及到的数组类元素个数不超过500个。参数`filter`中涉及到的`rules`数量不超过20个。参数`filter`的嵌套层级不超过3层。

#### filter
| 字段        | 类型     | 必选 | 描述        |
|-----------|--------|----|-----------|
| condition | string | 是  | 规则操作符     |
| rules     | array  | 是  | 过滤集群的范围规则 |

#### rules

过滤规则为三元组 `field`, `operator`, `value`

| 字段       | 类型     | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### time_condition

| 字段    | 类型     | 必选 | 描述           |
|-------|--------|----|--------------|
| oper  | string | 是  | 操作符，目前只支持and |
| rules | array  | 是  | 时间查询条件       |

#### time_condition.rules

| 字段    | 类型     | 必选 | 描述                          |
|-------|--------|----|-----------------------------|
| field | string | 是  | 取值为模型的字段名                   |
| start | string | 是  | 起始时间，格式为yyyy-MM-dd hh:mm:ss |
| end   | string | 是  | 结束时间，格式为yyyy-MM-dd hh:mm:ss |  

#### page

| 参数名称  | 参数类型   | 必选 | 描述           |
|-------|--------|----|--------------|
| start | int    | 是  | 记录开始位置       |
| limit | int    | 是  | 每页限制条数,最大200 |
| sort  | string | 否  | 排序字段         |

### 调用示例

```json
{
  "bk_biz_id": 2,
  "fields": [
    "bk_set_name"
  ],
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "bk_set_name",
        "operator": "equal",
        "value": "test"
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
  },
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "bk_set_name"
  }
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 1,
    "info": [
      {
        "bk_set_name": "test",
        "default": 0
      }
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
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data

| 参数名称  | 参数类型  | 描述                     |
|-------|-------|------------------------|
| count | int   | 数据数量                   |
| info  | array | 结果集，其中，所有字段均为集群定义的属性字段 |

#### info

| 参数名称                 | 参数类型   | 描述                         |
|----------------------|--------|----------------------------|
| bk_set_name          | string | 集群名称                       |
| default              | int    | 0-普通集群，1-内置模块集合，默认为0       |
| bk_biz_id            | int    | 业务id                       |
| bk_capacity          | int    | 设计容量                       |
| bk_parent_id         | int    | 父节点的ID                     |
| bk_set_id            | int    | 集群id                       |
| bk_service_status    | string | 服务状态:1/2(1:开放,2:关闭)        |
| bk_set_desc          | string | 集群描述                       |
| bk_set_env           | string | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
| create_time          | string | 创建时间                       |
| last_time            | string | 更新时间                       |
| bk_supplier_account  | string | 开发商账号                      |
| description          | string | 数据的描述信息                    |
| set_template_version | array  | 集群模板的当前版本                  |
| set_template_id      | int    | 集群模板ID                     |
| bk_created_at        | string | 创建时间                       |
| bk_updated_at        | string | 更新时间                       |
| bk_created_by        | string | 创建人                        |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
