### Functional description

Inquire service related information accord to host ID

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
| bk_host_id | array     | yes  | Host ID array, ID number can not exceed 500|
| bk_biz_id  | int           | no    | Business ID |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "bk_host_id": [
        3,
        4
    ]
}
```

### Return Result Example

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": [
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 59,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 61,
      "bk_set_id": 12,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 4,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    }
  ]
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

Data field Description:

| Name| Type| Description|
|---|---|---|
| bk_biz_id|  int| Business ID |
| bk_host_id|  int |Host ID|
| bk_module_id|  int| Module ID|
| bk_set_id|  int |set ID|
| bk_supplier_account|  string| Developer account|

