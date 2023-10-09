### Functional description

 Batch create generic model instances (v3.10.2+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Parameter      | Type   | Required| Description               |
| -------- | ------ | ---- | ------------------ |
| bk_obj_id | string |yes   | The model id used to create, allowing only instances of the generic model to be created   |
| details   |  array |yes   | The maximum number of instance contents to be created can not exceed 200, and the contents are the attribute information of the model instance|

#### details

| Parameter            | Type   | Required| Description           |
| --------------- | ------ | ---- | -------------- |
| bk_inst_name      |  string |yes   | Instance name   |
| bk_asset_id      |  string |yes| Fixed capital No.      |
| bk_sn | string |no| Equipment SN|
| bk_operator | string |no| Maintainer|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"bk_switch",
    "details":[
        {
            "bk_inst_name":"s1",
            "bk_asset_id":"test_001",
            "bk_sn":"00000001",
            "bk_operator":"admin"
        },
        {
            "bk_inst_name":"s2",
            "bk_asset_id":"test_002",
            "bk_sn":"00000002",
            "bk_operator":"admin"
        },
        {
            "bk_inst_name":"s3",
            "bk_asset_id":"test_003",
            "bk_sn":"00000003",
            "bk_operator":"admin"
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
            "1":1001,
            "2":1002
        },
        "error_msg":{
            "0":"duplicated instances exist, fields [bk_asset_id: test_001] duplicated"
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
| success_created | map |key is the index of the instance in the parameter details, and value is the id of the successfully created instance|
| error_msg       |  map |key is the index of the instance in the parameter details, and value is the failure information          |