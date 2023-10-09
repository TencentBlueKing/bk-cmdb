### Functional description

Query event subscription

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
| page                |  object     | no     | Paging parameter                    |
| condition           |  object     | no     | Query criteria                    |

#### page

| Field      | Type      | Required   | Description                |
|-----------|------------|--------|----------------------|
| start     |   int       | yes  | Record start position         |
| limit     |   int       | yes  | Limit bars per page, Max. 200|
| sort      |   string    | no     | Sort field             |

#### condition

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| subscription_name  |string      | yes | This is sample data only and needs to be set as a field for the query|

### Request Parameters Example

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

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data
| Field   | Type         | Description              |
|-------|--------------|------------------|
| count | int          | Number of records          |
| info  | array |Detailed list of event subscriptions|

#### info
| Field                 | Type      | Description                                       |
|----------------------|-----------|--------------------------------------------|
| subscription_id      |  int       | Subscription ID                                     |
| subscription_name    |  string    | Subscription name                                     |
| system_name          |  string    | System name                                   |
| callback_url         |  string    | Callback address                                   |
| confirm_mode         |  string    | Callback success confirmation mode, optional: httpstatus，regular|
| confirm_pattern      |  string    | Callback success flag                               |
| subscription_form    |  string    | Subscriptions, separated by ","                          |
| timeout              |  int       | Timeout in seconds                         |
| operator             |  int       | The person who last updated this piece of data                     |
| last_time            |  int       | Update time                                   |
| statistics.total     |  int       | Total push                                   |
| statistics.failure   |  int       | Number of push failures                                 |