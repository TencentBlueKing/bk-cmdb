### Functional description

Clone host properties

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field        | Type   | Required   | Description                       |
|-------------|---------|--------|-----------------------------|
| bk_org_ip   |  string  |yes     | Source host intranet ip   |
| bk_dst_ip   |  string  |yes     | Target host intranet ip|
| bk_org_id   |  int  |yes     | Source host ID    |
| bk_dst_id   |  int  |yes     | Destination host ID|
| bk_biz_id   |  int     | yes     | Business ID           |
| bk_cloud_id | int     | no     | Cloud area ID               |


Note: cloning by using host intranet IP and cloning by using host identity ID can only be used in one of the two methods, and can not be mixed.

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":2,
    "bk_org_ip":"127.0.0.1",
    "bk_dst_ip":"127.0.0.2",
    "bk_cloud_id":0
}
```
Or

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":2,
    "bk_org_id": 10,
    "bk_dst_id": 11,
    "bk_cloud_id":0
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
    "data": null
}
```

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |
