### Function Description

Unsubscribe from an event

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field           | Type | Required | Description     |
| --------------- | ---- | -------- | --------------- |
| subscription_id | int  | Yes      | Subscription ID |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "subscription_id": 1
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