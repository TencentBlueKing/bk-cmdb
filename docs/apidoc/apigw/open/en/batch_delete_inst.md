### Description

Batch delete object instances (Permission: Model instance deletion permission)

### Parameters

| Name     | Type  | Required | Description                |
|----------|-------|----------|----------------------------|
| inst_ids | array | Yes      | Collection of instance IDs |

### Request Example

```python
{
    "bk_obj_id": "bk_firewall",
    "delete": {
        "inst_ids": [
            46, 47
        ]
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
