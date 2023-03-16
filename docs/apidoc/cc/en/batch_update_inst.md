### Functional description

Bulk update object instances

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type       | Required   | Description                            |
|---------------------|-------------|--------|----------------------------------|
| bk_obj_id           |  string      | yes  | Model ID                           |
| update              |  array| yes     | Instance updated fields and values             |

#### update
| Field         | Type   | Required| Description                           |
|--------------|--------|-------|--------------------------------|
| datas        |  object |yes    | The value of the field for which the instance is updated           |
| inst_id      |  int    | yes | Indicates the specific instance that datas uses for the update   |

#### datas
| Field         | Type   | Required| Description                           |
|--------------|--------|-------|--------------------------------|
| bk_inst_name | string |no    | Instance name, or any other custom field|

**Datas is an object of map type, key is a field defined by the model corresponding to the instance, and value is the value of the field**


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"test",
    "update":[
        {
          "datas":{
            "bk_inst_name":"batch_update"
          },
          "inst_id":46
         }
        ]
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
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |
