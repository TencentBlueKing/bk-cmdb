### 描述

根据容器拓扑获取container信息 (版本：v3.12.5+，权限：业务访问)

### 输入参数
| 字段               | 类型    | 必选 | 描述                                                     |
|------------------|-------|----|--------------------------------------------------------|
| bk_biz_id        | int   | 是  | 业务ID                                                   |
| bk_cluster_id    | int   | 否  | cluster在cmdb中的唯一ID                                     |
| bk_namespace_id  | int   | 否  | namespace在cmdb中的唯一ID                                   |
| bk_workload_id   | int   | 否  | workload在cmdb中的唯一ID                                    |
| pod_filter       | object | 否  | 根据pod的查询条件                                             |
| container_filter | object | 否  | 根据container的查询条件                                       |
| container_fields | array | 是  | container属性列表，控制返回结果的container里有哪些字段，能够加速接口请求和减少网络流量传输 |
| pod_fields       | array | 是  | pod属性列表，控制返回结果的pod里有哪些字段，能够加速接口请求和减少网络流量传输             |
| page             | object | 是  | 分页信息                                                   |


#### pod_filter & container_filter

pod和container属性字段过滤规则，该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

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

组装规则可参考: <https://github.com/TencentBlueKing/bk-cmdb/blob/master/pkg/filter/README.md>

#### page 字段说明

| 字段           | 类型     | 必选  | 描述                                                                         |
|--------------|--------|-----|----------------------------------------------------------------------------|
| start        | int    | 是   | 记录开始位置                                                                     |
| limit        | int    | 是   | 每页限制条数，最大500                                                               |
| sort         | string | 否   | container排序字段                                                       |
| enable_count | bool   | 是   | 是否获取查询对象数量的标记。如果此标记为true那么表示此次请求是获取数量，此时其余字段必须为初始化值，start为0，limit为:0，sort为"" |

**注意：**

- `enable_count` 如果此标记为true，表示此次请求是获取数量。此时其余字段必须为初始化值，start为0,limit为:0, sort为""。
- 必须设置分页参数，一次最大查询数据不超过500个。

### 调用示例
```json
{
  "bk_biz_id": 1,
  "bk_cluster_id": 1,
  "bk_namespace_id": 1, 
  "bk_workload_id": 1,
  "pod_filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "equal",
        "value": "pod1"
      },
      {
        "field":"labels",
        "operator":"filter_object",
        "value":{
          "condition":"AND",
          "rules":[
            {
              "field":"location",
              "operator":"equal",
              "value":"sz"
            },
            {
              "field":"env",
              "operator":"equal",
              "value":"stage"
            }
          ]
        }
      }
    ]
  },
  "container_filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "id",
        "operator": "equal",
        "value": 1
      },
      {
        "field":"container_uid",
        "operator": "equal",
        "value": "uid"
      }
    ]
  },
  "container_fields": [
    "id",
    "name",
    "container_uid"
  ],
  "pod_fields": [
    "id",
    "labels",
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
           
### 响应示例
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
        "container": {
          "id": 1,
          "name": "container1",
          "container_uid": "uid"
        },
        "pod": {
          "id": 1,
          "labels": {
            "location": "sz",
            "env": "dev"
          },
          "name": "pod1"
        },
        "topo": {
          "bk_biz_id": 1,
          "bk_cluster_id": 1,
          "bk_namespace_id": 1,
          "bk_workload_id": 1,
          "workload_type": "deployment",
          "bk_host_id": 1
        }
      }
    ]
  }
}
```

### 响应参数说明
| 参数名称     | 参数类型   | 描述                           |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| data    | object | 请求返回的数据                           |

#### data
| 参数名称  | 参数类型  | 描述   |
|-------|-------|------|
| count | int   | 记录条数 |
| info  | array | 实际数据 |

#### data.info[x]
| 参数名称        | 参数类型   | 描述          |
|-------------|--------|-------------| 
| container   | object | container数据 | 
| pod         | object | pod数据       | 
| topo        | object | 拓扑数据        | 


#### data.info[x].container
| 参数名称        | 参数类型   | 描述          |
|---------------|--------------|------------------------------------------------------------|
| id    | int    | 在cmdb中的唯一ID   |
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

#### data.info[x].pod
| 参数名称        | 参数类型   | 描述          |
|---------------|--------------|------------------------|
| id    | int    | 在cmdb中的唯一ID   |
| name           | string       | 名称                                                                                                                                                         |
| priority       | int          | 优先级                                                                                                                                                        |
| labels         | string map   | 标签，key和value均是string，官方文档：http://kubernetes.io/docs/user-guide/labels                                                                                      |
| ip             | string       | 容器网络IP                                                                                                                                                     |
| ips            | object array | 容器网络IP数组，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#podip-v1-core                                                             |
| volumes        | object array | 使用的卷信息，官方文档：https://kubernetes.io/zh/docs/concepts/storage/volumes/ ，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volume-v1-core |
| qos_class      | enum         | 服务质量，官方文档：https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/quality-service-pod/                                                              |
| node_selectors | string map   | 节点标签选择器，key和value均是string，官方文档：https://kubernetes.io/zh/docs/concepts/scheduling-eviction/assign-pod-node/                                                 |
| tolerations    | object array | 容忍度，格式：https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#toleration-v1-core                                                             |
| operator       | string array | pod负责人                                                                                                                                                     |
| containers     | object array | 容器数据                                                                                                                                                       |

#### data.info[x].topo
| 参数名称            | 参数类型   | 描述                   |
|-----------------|--------------|----------------------|
| bk_biz_id       | int | 业务ID                 |
| bk_cluster_id   | int| cluster在cmdb中的唯一ID   |
| bk_namespace_id | int| namespace在cmdb中的唯一ID |
| bk_workload_id  | int | workload在cmdb中的唯一ID  |
| workload_type   | int | workload类型           |
| bk_host_id      | int | 主机id                 |