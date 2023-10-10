### Functional description

Bulk delete object instances

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type       | Required   | Description                            |
|---------------------|-------------|--------|----------------------------------|
| bk_obj_id           |  string      | yes  | Model ID|
| delete      |  object |yes    | Delete|

#### delete
| Field                | Type       | Required   | Description                            |
|---------------------|-------------|--------|----------------------------------|
| inst_ids |  array   | yes   | Instance ID set     |

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "test",
    "delete":{
    "inst_ids":[123]
    }
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

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |
