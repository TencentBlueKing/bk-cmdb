### Functional description

The service idle machine, that fault machine and the module to be recycle are obtained accord to the service ID

### Request Parameters

{{ common_args_desc }}


#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | yes  | Business ID |

### Request Parameters Example

```python

{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":0
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
    "bk_set_id": 2,
    "bk_set_name": "Idle machine",
    "module": [
      {
        "bk_module_id": 3,
        "bk_module_name": "Idle machine",
        "default": 1,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 4,
        "bk_module_name": "Faulty machine",
        "default": 2,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 5,
        "bk_module_name": "To be recycled",
        "default": 3,
        "host_apply_enabled": false
      }
    ]
  }
}
```

### Return Result Parameters Description
#### response
| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |


#### Data description
| Field      | Type      | Description      |
|-----------|------------|------------|
|bk_set_id | int64 |The instance ID of the set to which the idle machine, the failed machine, and the module to be recycled belong|
|bk_set_name | string |The instance name of the set to which the idle machine, the failed machine, and the module to be recycled belong|

#### Module description
| Field      | Type      | Description      |
|-----------|------------|------------|
|bk_module_id | int |The instance ID of the idle machine, failed machine, or module to be recycled|
|bk_module_name | string |The instance name of the idle machine, failed machine, or module to be recycled|
|default | int |Indicates the module type|
| host_apply_enabled| bool| Enable automatic application of host properties|
