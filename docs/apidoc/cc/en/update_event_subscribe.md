### Function Description

Modify subscription

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type   | Required | Description                                                  |
| ------------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_supplier_account | string | Yes      | Developer account                                            |
| subscription_id     | int    | Yes      | Subscription ID                                              |
| subscription_name   | string | Yes      | Name of the subscription                                     |
| system_name         | string | Yes      | Name of the system for subscription events                   |
| callback_url        | string | Yes      | Callback function                                            |
| confirm_mode        | string | Yes      | Event sending success verification mode, optional 1-httpstatus,2-regular |
| confirm_pattern     | string | Yes      | Callback httpstatus or regular expression                    |
| subscription_form   | string | Yes      | Events for subscription, separated by commas                 |
| timeout             | int    | Yes      | Timeout for sending events                                   |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_supplier_account": "0",
  "subscription_name":"mysubscribe",
  "subscription_id": 2,
  "system_name":"SystemName",
  "callback_url":"http://127.0.0.1:8080/callback",
  "confirm_mode":"httpstatus",
  "confirm_pattern":"200",
  "subscription_form":"hostcreate",
  "timeout":10
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