### Description

Lock hosts based on a list of host IDs. For newly added hosts, if the host has already been locked, it will also
indicate successful locking (Version: v3.8.6, Permission: Business host editing permission)

### Parameters

| Name    | Type      | Required | Description      |
|---------|-----------|----------|------------------|
| id_list | int array | Yes      | List of host IDs |

### Request Example

```python
{
   "id_list":[1, 2, 3]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null,
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| data       | object | Data returned by the request                                      |
| permission | object | Permission information                                            |
