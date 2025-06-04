### 描述

查询被引用的模型的实例列表(
版本：v3.10.30+，权限：传业务id时，代表是从业务的视角去查询引用的模型实例列表，当模型是业务时，鉴业务查看权限，否则鉴业务访问权限；其他的鉴对应模型实例的查看权限)

### 输入参数

| 参数名称           | 参数类型         | 必选 | 描述                                            |
|----------------|--------------|----|-----------------------------------------------|
| bk_biz_id      | string       | 否  | 业务id                                          |
| bk_obj_id      | string       | 是  | 源模型ID                                         |
| bk_property_id | string       | 是  | 源模型引用该模型的表格字段的唯一标识                            |
| filter         | object       | 否  | 被引用的模型实例的查询条件                                 |
| fields         | string array | 否  | 被引用的模型的属性列表，控制返回结果的实例里有哪些字段，能够加速接口请求和减少网络流量传输 |
| page           | object       | 是  | 分页信息                                          |

#### filter 字段说明

被引用的模型的属性字段过滤规则，用于根据属性字段搜索数据。该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

##### 组合过滤规则

由其它规则组合而成的过滤规则，组合的规则间支持逻辑与/或关系

| 参数名称      | 参数类型   | 必选 | 描述                              |
|-----------|--------|----|---------------------------------|
| condition | string | 是  | 组合查询条件，支持 `AND` 和 `OR` 两种方式     |
| rules     | array  | 是  | 查询规则，可以是 `组合过滤规则` 或 `原子过滤规则` 类型 |

##### 原子过滤规则

基础的过滤规则，表示对某一个字段进行过滤的规则。任何过滤规则都直接是原子过滤规则, 或由多个原子过滤规则组合而成

| 参数名称     | 参数类型                          | 必选 | 描述                                                                                                |
|----------|-------------------------------|----|---------------------------------------------------------------------------------------------------|
| field    | string                        | 是  | 被引用的模型的属性字段                                                                                       |
| operator | string                        | 是  | 操作符，可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | 不同的field和operator对应不同的value格式 | 否  | 操作值                                                                                               |

组装规则可参考: <https://github.com/TencentBlueKing/bk-cmdb/blob/v3.10.x/pkg/filter/README.md>

#### page 字段说明

| 参数名称         | 参数类型   | 必选 | 描述                                                                           |
|--------------|--------|----|------------------------------------------------------------------------------|
| start        | int    | 是  | 记录开始位置                                                                       |
| limit        | int    | 是  | 每页限制条数，最大500                                                                 |
| sort         | string | 否  | 排序字段                                                                         |
| enable_count | bool   | 是  | 是否获取查询对象数量的标记。如果此标记为true那么表示此次请求是获取数量，此时其余字段必须为初始化值，start为0，limit为:0，sort为"" |

### 调用示例

#### 获取详细信息请求参数示例

```json
{
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "test"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "operator",
            "operator": "not_in",
            "value": [
              "me"
            ]
          },
          {
            "field": "bk_inst_id",
            "operator": "equal",
            "value": 123
          }
        ]
      }
    ]
  },
  "fields": [
    "name",
    "description"
  ],
  "page": {
    "start": 0,
    "limit": 2,
    "sort": "name",
    "enable_count": false
  }
}
```

#### 获取数量请求参数示例

```json
{
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "filter": {
    "field": "name",
    "operator": "equal",
    "value": "test"
  },
  "page": {
    "enable_count": true
  }
}
```

### 响应示例

#### 获取详细信息返回结果示例

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
        "name": "test1",
        "description": "test instance 1"
      },
      {
        "name": "test2",
        "description": "test instance 2"
      }
    ]
  }
}
```

#### 获取数量返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 5,
    "info": []
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

| 参数名称  | 参数类型  | 描述                              |
|-------|-------|---------------------------------|
| count | int   | 记录条数                            |
| info  | array | 被引用的模型的实例的实际数据，仅返回fields里设置了的字段 |

#### info

| 参数名称        | 参数类型   | 描述                    |
|-------------|--------|-----------------------|
| name        | string | 名称，此处仅为示例，实际字段由模型属性决定 |
| description | string | 描述，此处仅为示例，实际字段由模型属性决定 |
