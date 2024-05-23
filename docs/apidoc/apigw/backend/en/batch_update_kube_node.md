### Description

Update Container Node Attribute Fields (Version: v3.12.1+, Permission: Edit Container Node Permission)

### Parameters

| Name      | Type   | Required | Description                         |
|-----------|--------|----------|-------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                         |
| ids       | array  | Yes      | List of node IDs to be updated      |
| data      | object | Yes      | Node attribute fields to be updated |

#### data

| Name              | Type        | Required | Description                                                           |
|-------------------|-------------|----------|-----------------------------------------------------------------------|
| labels            | json object | No       | Labels                                                                |
| taints            | string      | No       | Cluster name                                                          |
| unschedulable     | bool        | No       | Set whether it can be scheduled                                       |
| hostname          | string      | No       | Hostname                                                              |
| runtime_component | string      | No       | Runtime component                                                     |
| kube_proxy_mode   | string      | No       | Kube-proxy proxy mode                                                 |
| pod_cidr          | string      | No       | Allocation range of Pod addresses on this node, e.g., 172.17.0.128/26 |

**Note:**

- labels and taints need to be updated as a whole.
- The number of clusters to be updated at once should not exceed 100.

### Request Example

```json
{
  "bk_biz_id": 2,
  "ids": [
    1,
    2
  ],
  "data": {
    "labels": {
      "env": "test"
    },
    "taints": {
      "type": "gpu"
    },
    "unschedulable": false,
    "hostname": "xxx",
    "runtime_component": "runtime_component",
    "kube_proxy_mode": "ipvs",
    "pod_cidr": "127.0.0.1/26"
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
  "data": null
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
