### Functional description

create a new container node (v3.12.1+, permission: kube Node Creation Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type  | Required | Description                                                   |
|-----------|-------|----------|---------------------------------------------------------------|
| bk_biz_id | int   | yes      | business ID                                                   |
| data      | array | yes      | the specific information of the node that needs to be created |

#### data[x]

| Field             | Type   | Required | Description                                                                          |
|-------------------|--------|----------|--------------------------------------------------------------------------------------|
| bk_cluster_id     | int    | yes      | ID of the container cluster in cmdb                                                  |
| bk_host_id        | int    | yes      | associated host ID                                                                   |
| uid               | string | yes      | the own ID of the container cluster                                                  |
| name              | string | yes      | node name                                                                            |
| roles             | string | no       | node roles                                                                           |
| labels            | object | no       | label                                                                                |
| taints            | object | no       | taints                                                                               |
| unschedulable     | bool   | no       | Whether to turn off schedulable, true means not schedulable, false means schedulable |
| internal_ip       | array  | no       | internal ip                                                                          |
| external_ip       | array  | no       | external ip                                                                          |
| hostname          | string | no       | hostname                                                                             |
| runtime_component | string | no       | runtime components                                                                   |
| kube_proxy_mode   | string | no       | kube-proxy proxy mode                                                                |
| pod_cidr          | string | no       | The allocation range of the Pod address of this node                                 |

### Request Parameters Example

```json
 {
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
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
        "127.0.0.2"
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
- no more than 100 nodes can be created at one time.

### Return Result Example

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
  },
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

**Note:**

- The order of the node ID array in the returned data is consistent with the order of the array data in the parameter.

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                        |
|------------|--------|------------------------------------------------------------------------------------|
| result     | bool   | Whether the request succeeded or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                         |
| message    | string | Error message returned by request failure                                          |
| permission | object | Permission information                                                             |
| data       | object | Data returned by request                                                           |
| request_id | string | Request chain id                                                                   |

### data

| Name | Type  | Description                   |
|------|-------|-------------------------------|
| ids  | array | list of kube node IDs created |
