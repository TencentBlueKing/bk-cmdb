### Function Description

Delete the association between instances based on the ID of the instance relationship (Version: v3.5.40, Permission: Model instance deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                  |
| :-------- | :----- | :------- | :----------------------------------------------------------- |
| id        | int    | Yes      | ID of the instance relationship (Note: not the identity ID of the model instance), up to 500 |
| bk_obj_id | string | Yes      | The unique name of the source model of the relationship      |

### Request Parameter Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": 2
}
```

### Return Result Parameters Description

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | int    | Number of deleted associations                               |