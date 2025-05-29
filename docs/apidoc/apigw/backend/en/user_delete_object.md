### Description

Delete a model based on the model ID (Permission: Model deletion permission)

### Parameters

| Name | Type | Required | Description                         |
|------|------|----------|-------------------------------------|
| id   | int  | Yes      | ID of the data record to be deleted |

### Request Example

```python
{
    "id" : 0
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
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
