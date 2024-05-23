### Description

Delete a dynamic group (Version: v3.9.6, Permission: Dynamic group deletion permission)

### Parameters

| Name      | Type   | Required | Description                  |
|-----------|--------|----------|------------------------------|
| bk_biz_id | int    | Yes      | Business ID                  |
| id        | string | Yes      | Dynamic group primary key ID |

### Request Example

```json
{
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
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

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
