### 描述

查询业务集下容器信息 (版本：v3.14.8+，权限：业务集访问，该接口为bk-ops专用接口，可能会随时调整，其它服务请勿使用)

### 输入参数

| 字段               | 类型    | 必选 | 描述                                   |
|------------------|-------|----|--------------------------------------|
| bk_biz_set_id    | int   | 是  | 业务集ID                                |
| bk_container_ids | array | 是  | 容器ID列表，最多500条                        |
| container_fields | array | 是  | container属性列表，控制返回结果的container里有哪些字段 |
| pod_fields       | array | 是  | pod属性列表，控制返回结果的pod里有哪些字段             |
| host_fields      | array | 是  | host属性列表，控制返回结果的host里有哪些字段           |

### 调用示例

```json
{
  "bk_biz_set_id": 1,
  "bk_container_ids": [
    1,
    2
  ],
  "container_fields": [
    "id",
    "name",
    "container_uid"
  ],
  "pod_fields": [
    "id",
    "name"
  ],
  "host_fields": [
    "bk_host_id",
    "bk_agent_id",
    "bk_cloud_id"
  ]
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [
    {
      "container": {
        "id": 1,
        "name": "container1",
        "container_uid": "uid"
      },
      "pod": {
        "id": 1,
        "name": "pod1"
      },
      "host": {
        "bk_host_id": 1,
        "bk_agent_id": "xxxxxxxxxxx",
        "bk_cloud_id": 0
      }
    }
  ]
}
```

### 响应参数说明

| 参数名称       | 参数类型         | 描述                         |
|------------|--------------|----------------------------|
| result     | bool         | 请求成功与否。true:请求成功；false请求失败 |
| code       | int          | 错误编码。 0表示success，>0表示失败错误  |
| message    | string       | 请求失败返回的错误信息                |
| permission | object       | 权限信息                       |
| data       | object array | 请求返回的数据                    |

#### data[x]

| 参数名称      | 参数类型   | 描述          |
|-----------|--------|-------------| 
| container | object | container数据 | 
| pod       | object | pod数据       | 
| host      | object | host数据      |

#### data[x].container

| 参数名称          | 参数类型         | 描述                                                                                                                                                                                                        |
|---------------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id            | int          | 在cmdb中的唯一ID                                                                                                                                                                                               |
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

#### data[x].pod

| 参数名称           | 参数类型         | 描述                                                                                                                                                          |
|----------------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id             | int          | 在cmdb中的唯一ID                                                                                                                                                 |
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

#### data[x].host

| 参数名称        | 参数类型   | 描述                  |
|-------------|--------|---------------------|
| bk_host_id  | int    | 主机ID                |
| bk_agent_id | string | Agent ID            |
| bk_cloud_id | int    | 云区域ID               |
| 其他属性字段      | 对应属性类型 | 其他自定义属性字段，与模型属性定义对应 |
