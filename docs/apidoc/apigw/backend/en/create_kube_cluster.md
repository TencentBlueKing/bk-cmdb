### Description

Create Container Cluster (v3.12.1+, Permission: Create Permission for Container Cluster)

### Parameters

| Name              | Type   | Required | Description                                                                                                 |
|-------------------|--------|----------|-------------------------------------------------------------------------------------------------------------|
| bk_biz_id         | int    | Yes      | Business ID                                                                                                 |
| name              | string | Yes      | Cluster name                                                                                                |
| scheduling_engine | string | No       | Scheduling engine                                                                                           |
| uid               | string | Yes      | Cluster's own ID                                                                                            |
| xid               | string | No       | Associated cluster ID                                                                                       |
| version           | string | No       | Cluster version                                                                                             |
| network_type      | string | No       | Network type                                                                                                |
| region            | string | No       | Region                                                                                                      |
| vpc               | string | No       | VPC network                                                                                                 |
| network           | array  | No       | Cluster network                                                                                             |
| type              | string | Yes      | Cluster type. Enumeration values: INDEPENDENT_CLUSTER (Independent Cluster), SHARE_CLUSTER (Shared Cluster) |
| environment       | string | No       | Environment                                                                                                 |
| bk_project_id     | string | No       | Project ID                                                                                                  |
| bk_project_name   | string | No       | Project name                                                                                                |
| bk_project_code   | string | No       | Project English name                                                                                        |

### Request Example

```json
{
  "bk_biz_id": 2,
  "name": "cluster",
  "scheduling_engine": "k8s",
  "uid": "xxx",
  "xid": "xxx",
  "version": "1.1.0",
  "network_type": "underlay",
  "region": "xxx",
  "vpc": "xxx",
  "network": [
    "127.0.0.0/21"
  ],
  "type": "INDEPENDENT_CLUSTER",
  "environment": "xxx",
  "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
  "bk_project_name": "test",
  "bk_project_code": "test"
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
    "id": 1
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Request return data                                                         |
