### Description

Delete host locks based on a list of host IDs (v3.8.6, Permission: Business host editing permission)

### Parameters

| Name    | Type      | Required | Description      |
|---------|-----------|----------|------------------|
| id_list | int array | Yes      | List of host IDs |

### Request Example

```python
{
   "id_list": [1, 2, 3]
}
```

### Response Example

```python
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
