### Functional description

search subscription

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
| bk_biz_id           | string     | Yes     | Business ID                      |
| bk_supplier_account | string     | Yes     | Supplier account,please fill '0' by independent deployment  |
| page                | object     | Yes     | Page parameters                    |
| condition           | object     | No     | Search condition                    |
| fields              |string array| No     | Search fields                  |

##### page

| Field      |  Type      | Required	   |  Description                |
|-----------|------------|--------|----------------------|
| start     |  int       | Yes     | The record of start position         |
| limit     |  int       | Yes     | Limit number of each page,maximum 200 |
| sort      |  string    | No     | Sort fields             |

##### condition

| Field      |  Type      | Required	   |  Description      |
|-----------|------------|--------|------------|
| subscription_name  |string      |Yes      | Here is a sample data, which needs to be set as the query condition field |

### Request Parameters Example

```python
{
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

#### data
| Field | Type    | Description             |
| ----- | ------- | ----------------------- |
| count | integer | total count             |
| info  | array   | event subscription list |

#### info
| Field                 | Type      | Description                                       |
|----------------------|-----------|--------------------------------------------|
| subscription_id      | int       | Subscription ID                                     |
| subscription_name    | string    | Subscription name                                      |
| system_name          | string    | System name                                    |
| callback_url         | string    | Callback url                                    |
| confirm_mode         | string    | Confirm mode,optional: httpstatus,regular |
| confirm_pattern      | string    | Confirm pattern                                |
| subscription_form    | string    | Subscription form, separated by ','                          |
| timeout              | int       | Timeout, unit: second                         |
| operator             | int       | The last editor of data                     |
| last_time            | int       | Last update time                                     |
| statistics.total     | int       | Total statistics                                   |
| statistics.failure   | int       | Failure statistics                                 |
