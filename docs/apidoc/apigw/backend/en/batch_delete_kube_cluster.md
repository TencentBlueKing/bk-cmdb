### Description

Delete container cluster (v3.12.1+, Permission: Container cluster deletion permission)

### Parameters

| Name      | Type  | Required | Description                                  |
|-----------|-------|----------|----------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID of the container cluster         |
| ids       | array | Yes      | List of IDs of the container cluster in CMDB |

**Note:**

- Users need to ensure that there are no associated resources (such as namespace, pod, node workload, etc.) under the
  clusters to be deleted, otherwise, deletion will fail.
- The number of clusters to be deleted in one go should not exceed 10.

### Request Example

```json
{
  "bk_biz_id": 2,
  "ids": [
    1,
    2
  ]
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | No data returned                                                            |
