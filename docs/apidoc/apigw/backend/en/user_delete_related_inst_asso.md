### Description

Delete the association between instances based on the ID of the instance relationship (Version: v3.5.40, Permission:
Model instance deletion permission)

### Parameters

| Name      | Type   | Required | Description                                                                                  |
|-----------|--------|----------|----------------------------------------------------------------------------------------------|
| id        | int    | Yes      | ID of the instance relationship (Note: not the identity ID of the model instance), up to 500 |
| bk_obj_id | string | Yes      | The unique name of the source model of the relationship                                      |

### Request Example

```json
{
    "id":[1,2],
    "bk_obj_id": "abc"
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": 2
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | int    | Number of deleted associations                                     |
