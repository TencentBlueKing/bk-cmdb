### Description

Update Container Cluster Attribute Fields (Version: v3.12.1+, Permission: Edit Permission for Container Cluster)

### Parameters

| Name      | Type   | Required | Description                    |
|-----------|--------|----------|--------------------------------|
| bk_biz_id | int    | Yes      | Business ID                    |
| ids       | array  | Yes      | Unique IDs of clusters in cmdb |
| data      | object | Yes      | Data to be updated             |

#### data

| Name            | Type   | Required | Description          |
|-----------------|--------|----------|----------------------|
| name            | string | No       | Cluster name         |
| version         | string | No       | Cluster version      |
| network_type    | string | No       | Network type         |
| region          | string | No       | Region               |
| network         | array  | No       | Cluster network      |
| environment     | string | No       | Environment          |
| bk_project_id   | string | No       | Project ID           |
| bk_project_name | string | No       | Project name         |
| bk_project_code | string | No       | Project English name |

**Note:**

- The number of clusters to be updated at once should not exceed 100.
- This interface does not support updating the cluster type. If you need to update the cluster type, please use
  the `update_kube_cluster_type` interface.

### Request Example

```json
{
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
    "environment": "xxx",
    "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
    "bk_project_name": "test",
    "bk_project_code": "test"
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
