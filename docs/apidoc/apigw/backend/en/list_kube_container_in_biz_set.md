### Description

Query container information under a business set (version: v3.14.8+, permission: biz set access, this interface is
dedicated to bk-ops and may be adjusted at any time. Please do not use it for other services)

### Parameters

| Name             | Type         | Required | Description                    |
|------------------|--------------|----------|--------------------------------|
| bk_biz_set_id    | int64        | Yes      | Business set ID                |
| bk_container_ids | int64 array  | Yes      | Container ID list, maximum 500 |
| container_fields | string array | Yes      | Container attribute list       |
| pod_fields       | string array | Yes      | Pod attribute list             |
| host_fields      | string array | Yes      | Host attribute list            |

### Request Example

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

### Response Example

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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure              |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |

#### data[x]

| Name      | Type   | Description    |
|-----------|--------|----------------|
| container | object | Container data |
| pod       | object | Pod data       |
| host      | object | Host data      |

#### data[x].container

| Name          | Type         | Description                                                                                                                                                                                                               |
|---------------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id            | int          | Unique ID in cmdb                                                                                                                                                                                                         |
| name          | string       | Container name                                                                                                                                                                                                            |
| container_uid | string       | Container uid                                                                                                                                                                                                             |
| image         | string       | Container image                                                                                                                                                                                                           |
| ports         | object array | Container ports, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core                                                                                                       |
| args          | string array | Startup arguments                                                                                                                                                                                                         |
| started       | timestamp    | Start time                                                                                                                                                                                                                |
| limits        | object       | Resource limits, docs: https://kubernetes.io/docs/concepts/policy/resource-quotas/                                                                                                                                        |
| requests      | object       | Resource requests, docs: https://kubernetes.io/docs/concepts/policy/resource-quotas/                                                                                                                                      |
| liveness      | object       | Liveness probe, docs: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core   |
| environment   | object array | Environment variables, docs: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envvar-v1-core |
| mounts        | object array | Volume mounts, docs: https://kubernetes.io/docs/tasks/configure-pod-container/configure-volume-storage/, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumemount-v1-core                 |

#### data[x].pod

| Name           | Type         | Description                                                                                                                                                                      |
|----------------|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id             | int          | Unique ID in cmdb                                                                                                                                                                |
| name           | string       | Name                                                                                                                                                                             |
| priority       | int          | Priority                                                                                                                                                                         |
| labels         | string map   | Labels, both key and value are strings, docs: http://kubernetes.io/docs/user-guide/labels                                                                                        |
| ip             | string       | Container network IP                                                                                                                                                             |
| ips            | object array | Container network IP array, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#podip-v1-core                                                           |
| volumes        | object array | Volume information used, docs: https://kubernetes.io/docs/concepts/storage/volumes/, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volume-v1-core |
| qos_class      | enum         | Quality of Service, docs: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/                                                                          |
| node_selectors | string map   | Node label selector, both key and value are strings, docs: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/                                              |
| tolerations    | object array | Tolerations, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#toleration-v1-core                                                                     |
| operator       | string array | Pod operator                                                                                                                                                                     |
| containers     | object array | Container data                                                                                                                                                                   |

#### data[x].host

| Name         | Type           | Description                                           |
|--------------|----------------|-------------------------------------------------------|
| bk_host_id   | int            | Host ID                                               |
| bk_agent_id  | string         | Agent ID                                              |
| bk_cloud_id  | int            | Cloud area ID                                         |
| other fields | attribute type | other user-defined fields, defined by model attribute |
