### Functional description

The Association between instances of the query model.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description|
|----------------------|------------|--------|-----------------------------|
| condition | string map     | yes | Query criteria|


condition params

| Field                 | Type      | Required	   | Description|
|---------------------|------------|--------|-----------------------------|
| bk_asst_id           |  string     |  yes  | Association type unique id of the model|
| bk_obj_id           |  string     |  yes  | Source model id|
| bk_asst_id           |  string     |  yes  | Target model id|


### Request Parameters Example

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_asst_id": "belong",
        "bk_obj_id": "bk_switch",
        "bk_asst_obj_id": "bk_host"
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
           "id": 27,
           "bk_supplier_account": "0",
           "bk_obj_asst_id": "test1_belong_biz",
           "bk_obj_asst_name": "1",
           "bk_obj_id": "test1",
           "bk_asst_obj_id": "biz",
           "bk_asst_id": "belong",
           "mapping": "n:n",
           "on_delete": "none",
           "ispre": null
        }
    ]
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

#### data

| Field       | Type     | Description|
|------------|----------|--------------|
| id| int64| The identity id of the model Association|
| bk_obj_asst_id|  string| The unique ID of the model Association.|
| bk_obj_asst_name|  string| Alias for the Association. |
| bk_asst_id|  string| Association type id|
| bk_obj_id|  string| Source model id|
| bk_asst_obj_id|  string| Target model id|
| mapping|  string| The mapping relationship of the Association relationship instance between the source model and the target model can be one of the following [1: 1,1:n, n: n]|
| on_delete|  string| The action to delete an Association is one [of none, delete_src, delete_dest], "none" does nothing, "delete_src" deletes an instance of the source model, and "delete_dest" deletes an instance of the target model.|
| bk_supplier_account | string |Developer account number   |
| ispre               |  bool         | True: preset field,false: Non-built-in field                             |