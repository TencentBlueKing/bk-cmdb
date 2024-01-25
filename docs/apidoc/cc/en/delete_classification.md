### Function Description

Delete a model category by model category ID. (Permission: Model category deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type | Required | Description             |
| ----- | ---- | -------- | ----------------------- |
| id    | int  | Yes      | Category data record ID |

**Note** Can only delete an empty model category. Deletion will fail if the category has models.

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id": 13
}
```

### Response Example

#### Deletion Successful

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

#### Deletion Failed (Category has models)

```python
{
    "result": false,
    "code": 1101029,
    "data": null,
    "message": "There is a model under the category, not allowed to delete",
    "permission": null,
    "request_id": "8c6b89e7f0cb4fad836f55d50f81f2c6"
}
```

#### Response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |