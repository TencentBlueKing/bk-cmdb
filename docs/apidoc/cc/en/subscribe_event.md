### Function Description

Subscribe to events

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field             | Type   | Required | Description                                                  |
| ----------------- | ------ | -------- | ------------------------------------------------------------ |
| subscription_name | string | Yes      | Name of the subscription                                     |
| system_name       | string | Yes      | Name of the system for subscribed events                     |
| callback_url      | string | Yes      | Callback function                                            |
| confirm_mode      | string | Yes      | Event send success verification mode, optional: 1-httpstatus, 2-regular |
| confirm_pattern   | string | Yes      | callback's httpstatus or regular expression                  |
| subscription_form | string | Yes      | Subscribed events, separated by commas                       |
| timeout           | int    | Yes      | Timeout for sending events                                   |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "subscription_name":"mysubscribe",
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
    "data":{
        "subscription_id": 1
    }
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
| data       | object | Request returned data                                        |

#### data

| Field            | Type | Description                                      |
| --------------- | ---- | ------------------------------------------------ |
| subscription_id | int  | Subscription ID for the newly added subscription |