### Functional description

According to the service id, the set id and the module id, the host computer under the designated service set module is uploaded to the idle machine module of the service

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field          | Type      | Required     | Description    |
|---------------|------------|----------|----------|
| bk_biz_id     |  int        | yes    | Business ID |
| bk_set_id     |  int        | yes    | Set id |
| bk_module_id  | int        | yes    | Module id   |


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":10,
    "bk_module_id":58,
    "bk_set_id":1
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
    "data": "sucess"
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
