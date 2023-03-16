### Functional description

 Batch create general model instance Association (v3.10.2+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Parameter           | Type   | Required| Description                     |
| -------------- | ------ | ---- | ------------------------ |
| bk_obj_id      |  string |yes   | Source model id                 |
| bk_asst_obj_id | string |yes   | Target model model id           |
| bk_obj_asst_id | string |yes   | The unique id of the relationship between models|
| details        |  array  |yes   | The content of batch creation Association relationship can not exceed 200 relationships        |

#### details

| Parameter            | Type   | Required| Description           |
| --------------- | ------ | ---- | -------------- |
| bk_inst_id      |  int |yes   | Source model instance id   |
| bk_asst_inst_id | int |yes   | Target model instance id|

#### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"bk_switch",
    "bk_asst_obj_id":"host",
    "bk_obj_asst_id":"bk_switch_belong_host",
    "details":[
        {
            "bk_inst_id":11,
            "bk_asst_inst_id":21
        },
        {
            "bk_inst_id":12,
            "bk_asst_inst_id":22
        }
    ]
}
```

### Return Result Example

```json
{
    "result":true,
    "code":0,
    "message":"",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data":{
        "success_created":{
            "0":73
        },
        "error_msg":{
             "1":"the association inst is not exist"
        }
    }
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

| Field            | Type| Description                                                     |
| -------------- | ---- | -------------------------------------------------------- |
| success_created | map |key is the index of the instance Association in the parameter details array, and value is the id of the successfully created instance Association|
| error_msg       |  map |key is the index of the instance Association in the parameter details array, and value is the failure information          |