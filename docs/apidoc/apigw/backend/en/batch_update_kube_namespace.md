### Description

Batch Update Namespace (Version: v3.12.1+, Permission: Edit Namespace Permission)

### Parameters

| Name      | Type   | Required | Description                   |
|-----------|--------|----------|-------------------------------|
| bk_biz_id | int    | Yes      | Business ID                   |
| data      | object | Yes      | Contains fields to be updated |
| ids       | array  | Yes      | Unique IDs in cc              |

#### data

| Name            | Type  | Required | Description                                  |
|-----------------|-------|----------|----------------------------------------------|
| labels          | map   | No       | Labels                                       |
| resource_quotas | array | No       | Namespace CPU and memory requests and limits |

### Request Example

```json
{
  "bk_biz_id": 1,
  "ids": [
    1,
    2,
    3
  ],
  "data": {
    "labels": {
      "test": "test",
      "test2": "test2"
    },
    "resource_quotas": [
      {
        "hard": {
          "memory": "20000Gi",
          "pods": "100",
          "cpu": "10k"
        },
        "scope_selector": {
          "match_expressions": [
            {
              "values": [
                "high"
              ],
              "operator": "In",
              "scope_name": "PriorityClass"
            }
          ]
        }
      }
    ]
  }
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "data": null,
  "message": "success",
  "permission": null,
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
