### Functional description

Submit to the host to the module to be recovered of the service

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id     |   int     | yes  | Business ID |
|   bk_set_id   |   int     | And bk_module_id at least    | Set ID |
|     bk_module_id |  int     | And bk_set_id fill in at least one | Module ID |
| bk_host_id    |   array   | yes                                | Host ID|

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "bk_set_id": 1,
    "bk_module_id": 1,
    "bk_host_id": [
        9,
        10
    ]
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
    "data": null
}
```
### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |