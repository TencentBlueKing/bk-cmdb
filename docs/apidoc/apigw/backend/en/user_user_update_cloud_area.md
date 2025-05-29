### Description

Update control area (Permission: Control area editing permission)

### Parameters

| Name          | Type   | Required | Description       |
|---------------|--------|----------|-------------------|
| bk_cloud_id   | int    | Yes      | Control area ID   |
| bk_cloud_name | string | No       | Control area name |

### Request Example

```json
{
    "bk_cloud_id": 5,
    "bk_cloud_name": "Control Area 1"
}
```

### Response Example

```json
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
