### Description

Batch delete set by set ID under a specified business ID (Permission: Business topology deletion permission)

### Parameters

| Name     | Type      | Required | Description                        |
|----------|-----------|----------|------------------------------------|
| inst_ids | int array | Yes      | Array of Cluster IDs to be deleted |

### Request Example

```python
{
    "bk_biz_id": 0,
    "delete": {
        "inst_ids": [123]
    }
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": "success"
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
