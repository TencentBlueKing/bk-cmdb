### Description

Delete specified business's cluster templates based on business ID and cluster template ID list (Permission: Cluster
template deletion permission)

### Parameters

| Name             | Type  | Required | Description              |
|------------------|-------|----------|--------------------------|
| bk_biz_id        | int   | Yes      | Business ID              |
| set_template_ids | array | Yes      | Cluster template ID list |

### Request Example

```json
{
    "bk_biz_id": 20,
    "set_template_ids": [59]
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
