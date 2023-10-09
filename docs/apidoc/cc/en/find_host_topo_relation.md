### Functional description

Get host-to-topology relationship

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id|  int| yes |  Business ID |
| bk_set_ids| array |no| List of set IDs, up to 200|
| bk_module_ids| array |no| Module ID list, up to 500|
| bk_host_ids| array |no| Host ID list, up to 500|
| page|  object| yes | Paging information|

#### Page field Description

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
|start| int| no | Get data offset position|
|limit| int| yes | Limit on the number of data pieces in the past, 200 is recommended|

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "page":{
        "start":0,
        "limit":10
    },
    "bk_biz_id":2,
    "bk_set_ids": [1, 2],
    "bk_module_ids": [23, 24],
    "bk_host_ids": [25, 26]
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "data": {
        "count": 2,
        "data": [
            {
                "bk_biz_id": 2,
                "bk_host_id": 2,
                "bk_module_id": 2,
                "bk_set_id": 2,
                "bk_supplier_account": "0"
            },
            {
                "bk_biz_id": 1,
                "bk_host_id": 1,
                "bk_module_id": 1,
                "bk_set_id": 1,
                "bk_supplier_account": "0"
            }
        ],
        "page": {
            "limit": 10,
            "start": 0
        }
    },
    "message": "success",
    "permission": null,
    "request_id": "f5a6331d4bc2433587a63390c76ba7bf"
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

#### Data field Description:

| Name| Type| Description|
|---|---|---|
| count|  int| Number of records|
| data|  object array |Data details list of host and set, module and set under service|
| page|  object| Page|

#### Data.data field Description:
| Name| Type| Description|
|---|---|---|
| bk_biz_id | int |Service ID|
| bk_set_id | int |Set ID|
| bk_module_id | int |Module ID|
| bk_host_id | int |Host ID|
| bk_supplier_account | string |Developer account number|

#### Data.page field Description:
| Name| Type| Description|
|---|---|---|
|start| int| Data offset position|
|limit| int| Limit on number of past data pieces|
