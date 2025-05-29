### Description

Modify the business enable status based on the business ID and status value (Permission: Business archive permission)

### Parameters

| Name      | Type   | Required | Description                                  |
|-----------|--------|----------|----------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                  |
| flag      | string | Yes      | Enable status, either "disabled" or "enable" |

### Request Example

```python
{
    "bk_biz_id": "3",
    "flag": "enable"
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

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
