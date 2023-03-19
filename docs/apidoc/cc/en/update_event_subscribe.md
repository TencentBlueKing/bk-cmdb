### Functional description

Modify Subscription

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                   | Type    | Required   | Description                                            |
|------------------------|----------|--------|--------------------------------------------------|
| bk_supplier_account    |  string   | yes  | Developer account number                                       |
| subscription_id        |  int      | yes  | Subscription ID                                           |
| subscription_name      |  string   | yes  | The name of the subscription                                       |
| system_name            |  string   | yes  | The name of the system to which the event is subscribed                             |
| callback_url           |  string   | yes  | Callback function                                         |
| confirm_mode           |  string   | yes  | Event sending success verification mode, optional 1 HttpStatus, 2 regular|
| confirm_pattern        |  string   | yes  | HttpStatus or regular for callback                       |
| subscription_form      |  string   | yes  | Subscribed events, separated by commas                            |
| timeout                |  int      | yes      | Send event timeout                                 |


### Request Parameters Example

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

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |
