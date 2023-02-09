### Functional description

Obtain the attribute information of the module instances under the specified service in batches according to the service ID and module instance ID list, plus the module attribute list you want to return (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  |yes     | Business ID |
| bk_ids  |  array  |yes     | Module instance ID list, i.e. bk_module_id list, can be filled in up to 500|
| fields  |   array   | yes  | Module attribute list, which controls the fields in the module information that returns the result|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "bk_ids": [
        56,
        57,
        58,
        59,
        60
    ],
    "fields": [
        "bk_module_id",
        "bk_module_name",
        "create_time"
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_module_id": 60,
            "bk_module_name": "sm1",
            "create_time": "2020-05-15T22:15:51.725+08:00",
            "default": 0
        },
        {
            "bk_module_id": 59,
            "bk_module_name": "m1",
            "create_time": "2020-05-12T21:04:47.286+08:00",
            "default": 0
        },
        {
            "bk_module_id": 58,
            "bk_module_name": "recycle",
            "create_time": "2020-05-12T21:03:37.238+08:00",
            "default": 3
        },
        {
            "bk_module_id": 57,
            "bk_module_name": "fault",
            "create_time": "2020-05-12T21:03:37.183+08:00",
            "default": 2
        },
        {
            "bk_module_id": 56,
            "bk_module_name": "idle",
            "create_time": "2020-05-12T21:03:37.122+08:00",
            "default": 1
        }
    ]
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

#### Data description
| Field      | Type      | Description      |
|-----------|------------|------------|
|bk_module_id | int |Module id|
|bk_module_name | string |Module name|
|default | int |Indicates the module type|
|create_time | string |Settling time|