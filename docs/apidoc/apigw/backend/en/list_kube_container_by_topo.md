### Description

Get container information based on container topology (version: v3.13.3+, permission: business view)

### Parameters
| Name             | Type   | Required | Description                                                                                                                                                                    |
|------------------|--------|----|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id        | int    | yes | business id                                                                                                                                                                    |
| bk_kube_nodes    | array  | no | Container topology node information, array constraint is 200, note: the final result will take the concatenation of the hit elements in the array                              |
| pod_filter       | object | no | Query condition based on pod                                                                                                                                                   | 
| container_filter | object | no | based on container's query conditions                                                                                                                                          |
| container_fields | array  | yes | A list of container attributes that control what fields are in the container of the returned result, which can speed up interface requests and reduce network traffic transfer |
| pod_fields       | array  | yes | A list of pod attributes that control which fields in the returned pods can speed up interface requests and reduce network traffic transfers                                   |
| page             | object | yes | Paging information                                                                                                                                                             |

#### bk_kube_nodes[x]
| fields | type | mandatory | description |
|------------------|-------|----|------------------------------|
| kind | string | yes | Resource type of this node, currently supported types: cluster, namespace, deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods                          |                        
| id | int | yes | Unique identifier for container topology node information |


#### pod_filter & container_filter

This parameter supports the following two types of filter rules for the pod and container attribute fields, of which the combined filter rules can be nested, with a maximum of two levels of nesting. The supported filter rules are listed below:

#### combined filter rule

This filter rule type defines filter rules composed of other rules, the combined rules support logic and/or
relationships

| Field     | Type   | Required | Description                                                                |
|-----------|--------|----------|----------------------------------------------------------------------------|
| condition | string | yes      | query criteria, support `AND` and `OR`                                     |
| rules     | array  | yes      | query rules, can be of `combined filter rule` or `atomic filter rule` type |

##### atomic filter rule

This filter rule type defines basic filter rules, which represent rules for filtering a field. Any filter rule is either
directly an atomic filter rule, or a combination of multiple atomic filter rules

| Field    | Type                                                                 | Required | Description                                                                                                          |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | yes      | container's field                                                                                                    |
| operator | string                                                               | yes      | operator, optional values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | different fields and operators correspond to different value formats | yes      | operand                                                                                                              |

Assembly rules can refer to: <https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md>

#### page

| Field        | Type   | Required | Description                                                                                                                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | yes      | Record start position                                                                                                                                                                                              |
| limit        | int    | yes      | Limit per page, maximum 500                                                                                                                                                                                        |
| sort         | string | no       | Sort the field                                                                                                                                                                                                     |
| enable_count | bool   | yes      | The flag defining Whether to get the the number of query objects. If this flag is true, then the request is to get the quantity. The remaining fields must be initialized, start is 0, and limit is: 0, sort is "" |

**Note:**

- `enable_count`If this flag is true, this request is a get quantity. The remaining fields must be initialized, start is
  0, and limit is: 0, sort is "."
- Paging parameters must be set, and the maximum query data at one time does not exceed 500.

### Request Example
```json
{
  "bk_biz_id": 1,
  "bk_kube_nodes":[
    {
      "id": 1,
      "kind": "namespace"
    },
    {
      "id": 3,
      "kind": "workload"
    }
  ],
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

### Response Example
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

### Response Parameters
| parameter name | parameter type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | The request was successful or not. true:The request was successful; falseThe request failed |
| code | int | Error code. 0 means success, >0 means failure error | | message | string | request success or not. true: request succeeded; false request failed.
| message | string | The error message returned by the failed request |
| permission | object | Permission information | | data | object | Error code.
| data | object | The data returned by the request |

#### data
| parameter name | parameter type | description |
| ------- | ------- | ------|
| count | int | number of records |
| info | array | actual data |

#### data.info[x]
| parameter name | parameter type | description |
| -------------| --------| -------------| 
| container | object | container data | 
| pod | object | pod data | 
| topo | object | topology data | 

#### data.info[x].container
| 参数名称        | 参数类型   | 描述                                                                                                                                                                                                                                              |
|---------------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id    | int    | unique id in cmdb                                                                                                                                                                                                                               |
| name          | string       | container name                                                                                                                                                                                                                                  |
| container_uid | string       | container uid                                                                                                                                                                                                                                   |
| image         | string       | container image                                                                                                                                                                                                                                 |
| ports         | object array | container port information list, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core                                                                                                             |
| args          | string array | start arguments                                                                                                                                                                                                                                 |
| started       | timestamp    | start time                                                                                                                                                                                                                                      |
| limits        | object       | resource limits, official documentation: https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/                                                                                                                                         |
| requests      | object       | resource requests, official documentation: https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/                                                                                                                                       |
| liveness      | object       | liveness probe, official documentation: https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core   |
| environment   | object array | environment variables, official documentation: https://kubernetes.io/zh/docs/tasks/inject-data-application/define-environment-variable-container/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envvar-v1-core |
| mounts        | object array | volume mounts, official documentation: https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-volume-storage/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumemount-v1-core                 |

#### data.info[x].pod
| 参数名称        | 参数类型   | 描述          |
|---------------|--------------|------------------------|
| id    | int    | unique id in cmdb       |
| name           | string       | pod name                                                                                                                                                                                            |
| priority       | int          | pod priority                                                                                                                                                                                        |
| labels         | string map   | pod labels, key and value are all string, official documentation: http://kubernetes.io/docs/user-guide/labels                                                                                       |
| ip             | string       | pod ip                                                                                                                                                                                              |
| ips            | object array | pod ip list, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#podip-v1-core                                                                                             |
| volumes        | object array | pod volume info list, official documentation: https://kubernetes.io/zh/docs/concepts/storage/volumes/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volume-v1-core |
| qos_class      | enum         | quality of service class, official documentation: https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/quality-service-pod/                                                               |
| node_selectors | string map   | node label selectors, key and value are all string, official documentation: https://kubernetes.io/zh/docs/concepts/scheduling-eviction/assign-pod-node/                                             |
| tolerations    | object array | pod toleration list, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#toleration-v1-core                                                                                |
| operator       | string array | pod operator                                                                                                                                                                                        |
| containers     | object array | container information list                                                                                                                                                                          |

#### data.info[x].topo
| 参数名称            | 参数类型   | 描述                          |
|-----------------|--------------|-----------------------------|
| bk_biz_id       | int | business id                 |
| bk_cluster_id   | int| cluster unique id in cmdb   |
| bk_namespace_id | int| namespace unique id in cmdb |
| bk_workload_id  | int | workload unique id in cmdb  |
| workload_type   | int | workload type               |
| bk_host_id      | int | host id                      |

