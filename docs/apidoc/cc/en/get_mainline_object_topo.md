### Functional description

Get the business topology of the mainline model

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx"
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
  "data": [
    {
      "bk_obj_id": "biz",
      "bk_obj_name": "business",
      "bk_supplier_account": "0",
      "bk_next_obj": "set",
      "bk_next_name": "set",
      "bk_pre_obj_id": "",
      "bk_pre_obj_name": ""
    },
    {
      "bk_obj_id": "set",
      "bk_obj_name": "set",
      "bk_supplier_account": "0",
      "bk_next_obj": "module",
      "bk_next_name": "module",
      "bk_pre_obj_id": "biz",
      "bk_pre_obj_name": "business"
    },
    {
      "bk_obj_id": "module",
      "bk_obj_name": "module",
      "bk_supplier_account": "0",
      "bk_next_obj": "host",
      "bk_next_name": "host",
      "bk_pre_obj_id": "set",
      "bk_pre_obj_name": "set"
    },
    {
      "bk_obj_id": "host",
      "bk_obj_name": "host",
      "bk_supplier_account": "0",
      "bk_next_obj": "",
      "bk_next_name": "",
      "bk_pre_obj_id": "module",
      "bk_pre_obj_name": "module"
    }
  ]
}
```

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data
| Field      | Type      | Description      |
|-----------|------------|------------|
|bk_obj_id | string |The unique ID of the model|
|bk_obj_name | string |Model name|
|bk_supplier_account | string |Developer account name|
|bk_next_obj | string |The next model unique ID of the current model|
|bk_next_name | string |Next model name for the current model|
|bk_pre_obj_id | string |Unique ID of the previous model of the current model|
|bk_pre_obj_name | string |The name of the model preceding the current model|
