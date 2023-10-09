### Functional description

Query the instance Association relationship of the model.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description          |
|----------------------|------------|--------|-----------------------------|
| condition | object     | yes | Query criteria|
| bk_obj_id           |  string     |  yes  | Source model id(v3.10+)|


#### condition

| Field                 | Type      | Required	   | Description         |
|---------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           |  string     |  yes  | The unique id of the model Association|
| bk_asst_id           |  string     |  no    | Unique id of the Association type|
| bk_asst_obj_id           |  string     |  no    | Target model id|


### Request Parameters Example

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_asst_obj_id": ""
    },
    "bk_obj_id": "xxx"
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
    "data": [{
        "id": 481,
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
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

#### data

| Field       | Type     | Description         |
|------------|----------|--------------|
|id|int|the association's unique id|
| bk_obj_asst_id|  string| Automatically generated model association id.|
| bk_obj_id|  string| Association relationship source model id|
| bk_asst_obj_id|  string| Association relation target model id|
| bk_inst_id|  int| Source model instance id|
| bk_asst_inst_id|  int| Target model instance id|

