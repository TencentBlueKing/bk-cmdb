### 功能描述

查询Container列表 (版本：v3.12.1+，权限：业务访问)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型           | 必选  | 描述                                                     |
|-----------|--------------|-----|--------------------------------------------------------|
| bk_biz_id | int          | 是   | 业务ID                                                   |
| bk_pod_id | int          | 是   | 所属pod的ID                                               |
| filter    | object       | 否   | container的查询条件                                         |
| fields    | string array | 是   | container属性列表，控制返回结果的container里有哪些字段，能够加速接口请求和减少网络流量传输 |
| page      | object       | 是   | 分页信息                                                   |

#### filter 字段说明

container的属性字段过滤规则，用于根据container的属性字段搜索数据。该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

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

#### page 字段说明

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

#### 获取详细信息请求参数示例

```json
{
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
  "bk_biz_id": 4,
  "bk_pod_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "container1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "container_uid",
            "operator": "not_in",
            "value": [
              "xxxxxx"
            ]
          },
          {
            "field": "image",
            "operator": "equal",
            "value": "xxx"
          }
        ]
      }
    ]
  },
  "fields": [
    "name",
    "container_uid"
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
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
  "bk_biz_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "container1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "container_uid",
            "operator": "not_in",
            "value": [
              "xxxxxx"
            ]
          },
          {
            "field": "image",
            "operator": "equal",
            "value": "xxx"
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

### 返回结果示例

#### 获取详细信息返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": {
    "count": 0,
    "info": [
      {
        "name": "container2",
        "container_uid": "xxx"
      },
      {
        "name": "container3",
        "container_uid": "xxx"
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
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": {
    "count": 100,
    "info": []
  }
}
```

### 返回结果参数

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

| 字段    | 类型    | 描述                             |
|-------|-------|--------------------------------|
| count | int   | 记录条数                           |
| info  | array | container实际数据，仅返回fields里设置了的字段 |

#### info[x]

| 字段            | 类型           | 描述                                                                                                                                                                                                        |
|---------------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| name          | string       | 名称                                                                                                                                                                                                        |
| container_uid | string       | 容器uid                                                                                                                                                                                                     |
| image         | string       | 镜像信息                                                                                                                                                                                                      |
| ports         | object array | 容器端口，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core                                                                                                        |
| args          | string array | 启动参数                                                                                                                                                                                                      |
| started       | timestamp    | 启动时间                                                                                                                                                                                                      |
| limits        | object       | 资源限制，官方文档：https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/                                                                                                                                  |
| requests      | object       | 申请资源大小，官方文档：https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/                                                                                                                                |
| liveness      | object       | 存活探针，官方文档：https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ ，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core |
| environment   | object array | 环境变量，官方文档：https://kubernetes.io/zh/docs/tasks/inject-data-application/define-environment-variable-container/ ，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envvar-v1-core      |
| mounts        | object array | 挂载卷，官方文档：https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-volume-storage/ ，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumemount-v1-core               |
