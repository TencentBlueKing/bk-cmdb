### 描述

查询Pod列表 (版本：v3.12.1+，权限：业务访问)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                                         |
|-----------|--------|----|--------------------------------------------|
| bk_biz_id | int    | 是  | 业务ID                                       |
| filter    | object | 否  | pod的查询条件                                   |
| fields    | array  | 是  | pod属性列表，控制返回结果的pod里有哪些字段，能够加速接口请求和减少网络流量传输 |
| page      | object | 是  | 分页信息                                       |

#### filter 字段说明

pod的属性字段过滤规则，用于根据pod的属性字段搜索数据。该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

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
| field    | string                        | 是  | pod的字段                                                                                            |
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

#### 获取详细信息请求参数示例

```json
{
  "bk_biz_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "pod1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "priority",
            "operator": "not_in",
            "value": [
              2,
              6
            ]
          },
          {
            "field": "qos_class",
            "operator": "equal",
            "value": "Burstable"
          }
        ]
      }
    ]
  },
  "fields": [
    "name",
    "priority"
  ],
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "name",
    "enable_count": false
  }
}
```

#### 获取数量请求参数示例

```json
{
  "bk_biz_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "pod1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "priority",
            "operator": "not_in",
            "value": [
              2,
              6
            ]
          },
          {
            "field": "qos_class",
            "operator": "equal",
            "value": "Burstable"
          }
        ]
      }
    ]
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
        "name": "pod2",
        "priority": 1
      },
      {
        "name": "pod3",
        "priority": 5
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
    "count": 10,
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

| 参数名称  | 参数类型  | 描述                       |
|-------|-------|--------------------------|
| count | int   | 记录条数                     |
| info  | array | pod实际数据，仅返回fields里设置了的字段 |

#### info[x]

| 参数名称           | 参数类型         | 描述                                                                                                                                                          |
|----------------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| name           | string       | 名称                                                                                                                                                          |
| priority       | int          | 优先级                                                                                                                                                         |
| labels         | string map   | 标签，key和value均是string，官方文档：http://kubernetes.io/docs/user-guide/labels                                                                                       |
| ip             | string       | 容器网络IP                                                                                                                                                      |
| ips            | object array | 容器网络IP数组，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#podip-v1-core                                                              |
| volumes        | object array | 使用的卷信息，官方文档：https://kubernetes.io/zh/docs/concepts/storage/volumes/ ，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volume-v1-core |
| qos_class      | enum         | 服务质量，官方文档：https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/quality-service-pod/                                                               |
| node_selectors | string map   | 节点标签选择器，key和value均是string，官方文档：https://kubernetes.io/zh/docs/concepts/scheduling-eviction/assign-pod-node/                                                  |
| tolerations    | object array | 容忍度，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#toleration-v1-core                                                              |
| operator       | string array | pod负责人                                                                                                                                                      |
| containers     | object array | 容器数据                                                                                                                                                        |
