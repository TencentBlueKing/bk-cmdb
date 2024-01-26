### Description

Update Container Cluster Type (v3.12.1+, Permission: Container Cluster Editing Permission)

### Parameters

| Name      | Type   | Required | Description                                                                                                |
|-----------|--------|----------|------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                |
| id        | int    | Yes      | Unique ID list of clusters in CMDB                                                                         |
| type      | string | Yes      | Cluster type. Enumerated values: INDEPENDENT_CLUSTER (Independent Cluster), SHARE_CLUSTER (Shared Cluster) |

### Request Example

```json
{
  "bk_biz_id": 2,
  "id": 1,
  "type": "INDEPENDENT_CLUSTER"
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

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | No data returned                                                    |
