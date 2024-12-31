### Description

Create Container Nodes (v3.12.1+, Permission: Container nodes creation permission)

### Parameters

| Name      | Type  | Required | Description                    |
|-----------|-------|----------|--------------------------------|
| bk_biz_id | int   | Yes      | Business ID                    |
| data      | array | Yes      | Node information to be created |

#### data[x]

| Name              | Type   | Required | Description                                                                        |
|-------------------|--------|----------|------------------------------------------------------------------------------------|
| bk_cluster_id     | int    | Yes      | Unique ID of the container cluster in CMDB                                         |
| bk_host_id        | int    | Yes      | Associated host ID                                                                 |
| name              | string | Yes      | Node name                                                                          |
| roles             | string | No       | Node type                                                                          |
| labels            | object | No       | Labels                                                                             |
| taints            | object | No       | Taints                                                                             |
| unschedulable     | bool   | No       | Whether to disable scheduling, true means not schedulable, false means schedulable |
| internal_ip       | array  | No       | Internal IP addresses                                                              |
| external_ip       | array  | No       | External IP addresses                                                              |
| hostname          | string | No       | Hostname                                                                           |
| runtime_component | string | No       | Runtime component                                                                  |
| kube_proxy_mode   | string | No       | Kube-proxy proxy mode                                                              |
| pod_cidr          | string | No       | Allocation range of Pod addresses for this node                                    |

### Request Example

```json
{
  "bk_biz_id": 2,
  "data": [
    {
      "bk_host_id": 1,
      "bk_cluster_id": 1,
      "name": "k8s",
      "roles": "master",
      "labels": {
        "env": "test"
      },
      "taints": {
        "type": "gpu"
      },
      "unschedulable": false,
      "internal_ip": [
        "127.0.0.1"
      ],
      "external_ip": [
        "127.0.0.1"
      ],
      "hostname": "xxx",
      "runtime_component": "runtime_component",
      "kube_proxy_mode": "ipvs",
      "pod_cidr": "127.0.0.1/26"
    },
    {
      "bk_host_id": 2,
      "bk_cluster_id": 1,
      "name": "k8s-node",
      "roles": "master",
      "labels": {
        "env": "test"
      },
      "taints": {
        "type": "gpu"
      },
      "unschedulable": false,
      "internal_ip": [
        "127.0.0.1"
      ],
      "external_ip": [
        "127.0.0.2"
      ],
      "hostname": "xxx",
      "runtime_component": "runtime_component",
      "kube_proxy_mode": "ipvs",
      "pod_cidr": "127.0.0.1/26"
    }
  ]
}
```

**Note:**

- internal_ip and external_ip cannot be empty at the same time.
- The number of nodes created at once should not exceed 100.

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "ids": [
      1,
      2
    ]
  }
}
```

**Note:**

- The order of the node ID array returned in the data field corresponds to the order of the array data in the
  parameters.

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned by the request                                      |
