### Functional description

batch update container cluster attribute fields (v3.12.1+, permission: kube cluster editing permissions)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                       |
|-----------|--------|----------|-----------------------------------|
| bk_biz_id | int    | yes      | business ID                       |
| ids       | array  | yes      | unique IDs of the cluster in cmdb |
| data      | object | yes      | data that needs to be updated     |

#### data

| Field           | Type   | Required | Description                             |
|-----------------|--------|----------|-----------------------------------------|
| name            | string | no       | cluster name                            |
| version         | string | no       | cluster version                         |
| network_type    | string | no       | network type                            |
| region          | string | no       | the region where the cluster is located |
| vpc             | string | no       | vpc network                             |
| network         | array  | no       | cluster network                         |
| bk_project_id   | string | no       | project_id                              |
| bk_project_name | string | no       | project name                            |
| bk_project_code | string | no       | project english name                    |

**Note:**

- the number of clusters to be updated at one time does not exceed 100
- this api does not support updating cluster type, please use the `update_kube_cluster_type` api to update it.

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "ids": [
    1
  ],
  "data": {
    "name": "cluster",
    "version": "1.20.6",
    "network_type": "underlay",
    "region": "xxx",
    "network": [
      "127.0.0.0/21"
    ],
    "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
    "bk_project_name": "test",
    "bk_project_code": "test"
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

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | object | Data returned by request                                                                |
