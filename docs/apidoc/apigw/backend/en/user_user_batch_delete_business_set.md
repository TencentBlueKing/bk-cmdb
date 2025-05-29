### Description

Delete business set (Version: v3.10.12+, Permission: business set deletion permission)

### Parameters

| Name           | Type  | Required | Description              |
|----------------|-------|----------|--------------------------|
| bk_biz_set_ids | array | Yes      | List of business set IDs |

### Request Example

```python
{
    "bk_biz_set_ids": [
        10,
        12
    ]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {},
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
