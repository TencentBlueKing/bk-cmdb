### Description

Batch create instances of a common model (Version: v3.10.2+, Permission: New instance permission)

### Parameters

| Name      | Type   | Required | Description                                                                                                         |
|-----------|--------|----------|---------------------------------------------------------------------------------------------------------------------|
| bk_obj_id | string | Yes      | Model ID for creation, only allows creating instances of common models                                              |
| details   | array  | Yes      | Content of instances to be created, up to 200 instances, content is the attribute information of the model instance |

#### details

| Name         | Type   | Required | Description           |
|--------------|--------|----------|-----------------------|
| bk_inst_name | string | Yes      | Instance name         |
| bk_asset_id  | string | Yes      | Asset ID              |
| bk_sn        | string | No       | Device SN             |
| bk_operator  | string | No       | Maintenance personnel |

### Request Example

```json
{
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

### Response Example

```json
{
    "result":true,
    "code":0,
    "message":"",
    "permission": null,
    "data":{
        "success_created":{
            "1":1001,
            "2":1002
        },
        "error_msg":{
            "0":"数据唯一性校验失败， [bk_asset_id: test_001] 重复"
        }
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned by the request                                      |

#### data

| Name            | Type | Description                                                                              |
|-----------------|------|------------------------------------------------------------------------------------------|
| success_created | map  | Key is the index in the details parameter, value is the instance ID created successfully |
| error_msg       | map  | Key is the index in the details parameter, value is the failure information              |
