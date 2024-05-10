### 描述

查询业务(权限：业务查询权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述                                                                 |
|---------------------|--------|----|--------------------------------------------------------------------|
| bk_supplier_account | string | 否  | 开发商账号                                                              |
| fields              | array  | 否  | 指定查询的字段，参数为业务的任意属性，如果不填写字段信息，系统会返回业务的所有字段                          |
| condition           | dict   | 否  | 查询条件，参数为业务的任意属性，如果不写代表搜索全部数据，(历史遗留字段，请勿继续使用，请用biz_property_filter) |
| biz_property_filter | object | 否  | 业务属性组合查询条件                                                         |
| time_condition | object | 否 | 按时间查询业务的查询条件 |
| page                | dict   | 否  | 分页条件                                                               |

Note: 业务分为两类，未归档的业务和已归档的业务。

- 若要查询已归档的业务，请在condition中增加条件`bk_data_status:disabled`。
- 若要查询未归档的业务，请不要带字段"bk_data_status",或者在condition中增条件`bk_data_status: {"$ne":disabled"}`。
- `biz_property_filter`与`condition`两个参数只能有一个生效，参数`condition`不建议继续使用。
- 参数`biz_property_filter` 中涉及到的数组类元素个数不超过500个。参数`biz_property_filter`中涉及到的`rules`
  数量不超过20个。参数`biz_property_filter`
  的嵌套层级不超过3层。

#### biz_property_filter

| 参数名称      | 参数类型   | 必选 | 描述   |
|-----------|--------|----|------|
| condition | string | 是  | 聚合条件 |
| rules     | array  | 是  | 规则   |

#### rules

| 参数名称     | 参数类型   | 必选 | 描述  |
|----------|--------|----|-----|
| field    | string | 是  | 字段  |
| operator | string | 是  | 操作符 |
| value    | object | 是  | 值   |

#### time_condition

| 字段   | 类型   | 必选 |  描述              |
|-------|--------|-----|--------------------|
| oper  | string | 是  | 操作符，目前只支持and |
| rules | array  | 是  | 时间查询条件         |

#### page

| 参数名称  | 参数类型   | 必选 | 描述                                                        |
|-------|--------|----|-----------------------------------------------------------|
| start | int    | 是  | 记录开始位置                                                    |
| limit | int    | 是  | 每页限制条数,最大200                                              |
| sort  | string | 否  | 排序字段，通过在字段前面增加 -，如 sort:&#34;-field&#34; 可以表示按照字段 field降序 |

### 调用示例

```json
{
  "fields": [
    "bk_biz_id",
    "bk_biz_name"
  ],
  "biz_property_filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "bk_biz_maintainer",
        "operator": "equal",
        "value": "admin"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "bk_biz_name",
            "operator": "in",
            "value": [
              "test"
            ]
          },
          {
            "field": "bk_biz_id",
            "operator": "equal",
            "value": 1
          }
        ]
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
    "sort": ""
  }
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": {
    "count": 1,
    "info": [
      {
        "bk_biz_id": 1,
        "bk_biz_name": "esb-test",
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

| 参数名称  | 参数类型  | 描述     |
|-------|-------|--------|
| count | int   | 记录条数   |
| info  | array | 业务实际数据 |

#### info

| 参数名称                | 参数类型   | 描述                   |
|---------------------|--------|----------------------|
| bk_biz_id           | int    | 业务id                 |
| bk_biz_name         | string | 业务名                  |
| bk_biz_maintainer   | string | 运维人员                 |
| bk_biz_productor    | string | 产品人员                 |
| bk_biz_developer    | string | 开发人员                 |
| bk_biz_tester       | string | 测试人员                 |
| time_zone           | string | 时区                   |
| language            | string | 语言, "1"代表中文, "2"代表英文 |
| bk_supplier_account | string | 开发商账号                |
| create_time         | string | 创建时间                 |
| last_time           | string | 更新时间                 |
| default             | int    | 表示业务类型               |
| operator            | string | 主要维护人                |
| life_cycle          | string | 业务生命周期               |
| bk_created_at       | string | 创建时间                 |
| bk_updated_at       | string | 更新时间                 |
| bk_created_by       | string | 创建人                  |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
