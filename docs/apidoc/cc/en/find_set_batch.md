### Functional description

Obtain the attribute details of the  set  under the specified service in batches according to the service id,  set  instance id list and the attribute list you want to obtain (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  |yes     | Business ID |
| bk_ids  | array  |yes     | Set instance ID list, namely bk_set_id list, can be filled with up to 500 |
| fields  |  array   | yes  | Set attribute list, which controls which fields are in the set information of the returned result |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "bk_ids": [
        11,
        12
    ],
    "fields": [
        "bk_set_id",
        "bk_set_name",
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
            "bk_set_id": 12,
            "bk_set_name": "ss1",
            "create_time": "2020-05-15T22:15:51.67+08:00",
            "default": 0
        },
        {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "create_time": "2020-05-12T21:04:36.644+08:00",
            "default": 0
        }
    ]
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
