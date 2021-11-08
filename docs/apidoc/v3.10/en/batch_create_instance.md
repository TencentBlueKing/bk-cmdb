### Functional description

batch create common model instances(v3.10.2+)


can not create mainline business model instance with this api.

### Request Parameters

{{ common_args_desc }}

#### Interface parameters

| Field     | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| bk_obj_id | string | Yes      | object id，only common instances allowed to be created       |
| details   | array  | Yes      | instance content to be created, maximum number is 200. the content is the attribute of the instance, please refer to the request parameters example |

### Request Parameters Example

```json
{
    "bk_obj_id":"bk_switch",
    "details":[
        {
            "bk_inst_name":"s1",
            "bk_asset_id":"test_001",
            "bk_sn":"00000001",
            "bk_operator":"admin",
            ...
        },
        {
            "bk_inst_name":"s2",
            "bk_asset_id":"test_002",
            "bk_sn":"00000002",
            "bk_operator":"admin",
            ...
        },
        {
            "bk_inst_name":"s3",
            "bk_asset_id":"test_003",
            "bk_sn":"00000003",
            "bk_operator":"admin",
            ...
        }
    ]
}
```

### Return Result Example

```json
{
    "result":true,
    "code":0,
    "message":"success",
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

#### data

| Field           | Type   | Description                    |
| --------------- | ------ | ------------------------------ |
| success_created | map | key is index of request's detail array, value is created instance id |
| error_msg       | map | key is index of request's detail array, value is error message       |