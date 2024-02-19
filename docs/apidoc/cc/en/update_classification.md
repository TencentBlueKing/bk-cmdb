### Function Description

Update model classification (Permission: Model group editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                  | Type   | Required | Description                                                  |
| ---------------------- | ------ | -------- | ------------------------------------------------------------ |
| id                     | int    | No       | The record ID of the target data, used as a condition for the update operation |
| bk_classification_name | string | No       | Classification name                                          |
| bk_classification_icon | string | No       | Model classification icon, the value can refer to [(classIcon.json)](https://chat.openai.com/c/resource_define/classIcon.json) |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id": 1,
    "bk_classification_name": "cc_test_new",
    "bk_classification_icon": "icon-cc-business"
}
```

### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Data returned by the request                                 |