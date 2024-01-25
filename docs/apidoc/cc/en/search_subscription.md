### Function Description

Query event subscription

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description           |
| --------- | ------ | -------- | --------------------- |
| page      | object | No       | Pagination parameters |
| condition | object | No       | Query conditions      |

#### page

| Field | Type   | Required | Description                                |
| ----- | ------ | -------- | ------------------------------------------ |
| start | int    | Yes      | Record start position                      |
| limit | int    | Yes      | Number of records per page, maximum is 200 |
| sort  | string | No       | Sorting field                              |

#### condition

| Field             | Type   | Required | Description                                                  |
| ----------------- | ------ | -------- | ------------------------------------------------------------ |
| subscription_name | string | Yes      | Subscription name (this is just an example data, it should be set to the field to be queried) |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account":"0",
    "condition":{
        "subscription_name":"name"
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"HostName"
    }
}
```

### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 1,
        "info": [
            {
                "subscription_id": 1,
                "subscription_name": "mysubscribe",
                "system_name": "SystemName",
                "callback_url": "http://127.0.0.1:8080/callback",
                "confirm_mode": "httpstatus",
                "confirm_pattern": "200",
                "time_out": 10,
                "subscription_form": "hostcreate",
                "operator": "user",
                "bk_supplier_account": "0",
                "last_time": "2017-09-19 16:57:07",
                "statistics": {
                    "total": 30,
                    "failure": 2
                }
            }
        ]
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

| Field  | Type  | Description                         |
| ----- | ----- | ----------------------------------- |
| count | int   | Number of records                   |
| info  | array | Details list of event subscriptions |

#### info

| Field               | Type   | Description                                                  |
| ------------------ | ------ | ------------------------------------------------------------ |
| subscription_id    | int    | Subscription ID                                              |
| subscription_name  | string | Subscription name                                            |
| system_name        | string | System name                                                  |
| callback_url       | string | Callback URL                                                 |
| confirm_mode       | string | Callback success confirmation mode, optional: httpstatus, regular |
| confirm_pattern    | string | Callback success flag                                        |
| subscription_form  | string | Subscription form, separated by ","                          |
| timeout            | int    | Timeout, unit: seconds                                       |
| operator           | int    | Last updated by                                              |
| last_time          | int    | Update time                                                  |
| statistics.total   | int    | Total number of pushes                                       |
| statistics.failure | int    | Number of failed pushes                                      |