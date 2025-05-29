### Description

Delete a control area based on the control area ID. (Permission: Control area deletion permission)

### Parameters

| Name        | Type | Required | Description     |
|-------------|------|----------|-----------------|
| bk_cloud_id | int  | Yes      | Control area ID |

### Request Example

```json
{
    "bk_cloud_id": 5
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": "success"
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
