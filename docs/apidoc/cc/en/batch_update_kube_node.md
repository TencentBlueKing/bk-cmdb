### Functional description

batch update container node attribute field (v3.12.1+, permission: kube Node Edit Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                   |
|-----------|--------|----------|-----------------------------------------------|
| bk_biz_id | int    | yes      | business ID                                   |
| ids       | array  | yes      | IDs of the node in cmdb                       |
| data      | object | yes      | Node attribute fields that need to be updated |

#### data

| Field             | Type        | Required | Description                                                                        |
|-------------------|-------------|----------|------------------------------------------------------------------------------------|
| labels            | json object | no       | label                                                                              |
| taints            | string      | no       | cluster name                                                                       |
| unschedulable     | bool        | no       | set whether to schedule                                                            |
| hostname          | string      | no       | host name                                                                          |
| runtime_component | string      | no       | runtime components                                                                 |
| kube_proxy_mode   | string      | no       | Kube-proxy proxy mode                                                              |
| pod_cidr          | string      | no       | The allocation range of the Pod address of this node, for example: 172.17.0.128/26 |

**注意：**

- Among them, labels and taints need to be updated as a whole.
- data field cannot be empty.
- The number of clusters to be updated at one time does not exceed 100.

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
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

### Return Result Example

```json
 {
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": null
}
```

### Return Result Parameters Description

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | object | Data returned by request                                                                |
