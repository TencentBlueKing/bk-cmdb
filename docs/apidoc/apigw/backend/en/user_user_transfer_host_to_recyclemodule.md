### Description

Submit hosts to the business's pending recycle module (Permission: Service Instance Edit Permission)

### Parameters

| Name         | Type  | Required                                                  | Description |
|--------------|-------|-----------------------------------------------------------|-------------|
| bk_biz_id    | int   | Yes                                                       | Business ID |
| bk_set_id    | int   | At least one of bk_set_id and bk_module_id must be filled | Cluster ID  |
| bk_module_id | int   | At least one of bk_set_id and bk_module_id must be filled | Module ID   |
| bk_host_id   | array | Yes                                                       | Host ID     |

### Request Example

```python
{
    "bk_biz_id": 1,
    "bk_set_id": 1,
    "bk_module_id": 1,
    "bk_host_id": [
        9,
        10
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
    "data": null
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
