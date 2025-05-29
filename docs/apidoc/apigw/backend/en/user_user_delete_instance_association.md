### Description

Delete the relationship between model instances based on the unique identity ID of the model instance relationship. (
Permission: Editing permission of source model instance and target model instance)

### Parameters

| Name      | Type   | Required | Description                                                           |
|-----------|--------|----------|-----------------------------------------------------------------------|
| id        | int    | Yes      | Unique identity ID of the model instance relationship                 |
| bk_obj_id | string | Yes      | Source or target model ID of the model instance relationship (v3.10+) |

### Request Example

```json
{
    "bk_obj_id": "test",
    "id": 1
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
