### Description

Query host locks based on host ID list (Version: v3.8.6, Permission: Business host edit permission)

### Parameters

| Name    | Type  | Required | Description  |
|---------|-------|----------|--------------|
| id_list | array | Yes      | Host ID list |

### Request Example

```python
{
   "id_list":[1, 2]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        1: true,
        2: false
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name | Type   | Description                                                                              |
|------|--------|------------------------------------------------------------------------------------------|
| data | object | Data returned by the request, where the key is ID, and the value is whether it is locked |
