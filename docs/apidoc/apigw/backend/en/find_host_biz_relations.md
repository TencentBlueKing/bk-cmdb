### Description

Search for business-related information based on host ID.

### Parameters

| Name       | Type  | Required | Description                                            |
|------------|-------|----------|--------------------------------------------------------|
| bk_host_id | array | Yes      | Array of host IDs, the number of IDs cannot exceed 500 |
| bk_biz_id  | int   | No       | Business ID                                            |

### Request Example

```json
{
    "bk_biz_id": 1,
    "bk_host_id": [
        3,
        4
    ]
}
```

### Response Example

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission": null,
  "data": [
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 59,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 61,
      "bk_set_id": 12,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 4,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    }
  ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

Explanation of the data field:

| Name                | Type   | Description       |
|---------------------|--------|-------------------|
| bk_biz_id           | int    | Business ID       |
| bk_host_id          | int    | Host ID           |
| bk_module_id        | int    | Module ID         |
| bk_set_id           | int    | Cluster ID        |
| bk_supplier_account | string | Developer account |
