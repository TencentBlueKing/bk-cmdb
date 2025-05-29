### Description

General Model Instance Relationship Query (Version: v3.10.1+, Permission: Model Instance Query Permission)

### Parameters

| Name       | Type   | Required | Description                                                                                                                                                                                                           |
|------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id  | int    | No       | Business ID, required for mainline model query                                                                                                                                                                        |
| bk_obj_id  | string | Yes      | Model ID                                                                                                                                                                                                              |
| conditions | object | No       | Combined query conditions, supports both AND and OR, can be nested, supports up to 3 layers, each OR condition supports up to 20, not specifying this parameter means matching all (i.e., conditions is null)         |
| fields     | array  | No       | Specify the fields to be returned. Fields that are not available will be ignored. If not specified, all fields will be returned (returning all fields will affect performance, it is recommended to return as needed) |
| page       | object | Yes      | Pagination settings                                                                                                                                                                                                   |

#### conditions

| Name      | Type   | Required | Description                                     |
|-----------|--------|----------|-------------------------------------------------|
| condition | string | Yes      | Rule operator                                   |
| rules     | array  | Yes      | Range condition rules for the selected business |

#### conditions.rules

| Name     | Type   | Required | Description                                                                                                                                    |
|----------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Condition field, optional values are id, bk_inst_id, bk_obj_id, bk_asst_inst_id, bk_asst_obj_id, bk_obj_asst_id, bk_asst_id                    |
| operator | string | Yes      | Operator, optional values are equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between, etc.                  |
| value    | -      | No       | Expected value of the condition field, different operators correspond to different value formats, array type values support up to 500 elements |

Detailed assembly rules can be referred
to: [bk-cmdb query builder](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### page

| Name  | Type   | Required | Description                                                                                              |
|-------|--------|----------|----------------------------------------------------------------------------------------------------------|
| start | int    | Yes      | Record start position                                                                                    |
| limit | int    | Yes      | Number of records per page, default is 500                                                               |
| sort  | string | No       | Retrieval sorting, follow the MongoDB semantic format {KEY}:{ORDER}, default sorting is by creation time |

### Request Example

```json
{
    "bk_obj_id":"bk_switch",
    "conditions":{
        "condition": "AND",
        "rules": [
            {
                "field": "bk_obj_asst_id",
                "operator": "equal",
                "value": "bk_switch_connect_host"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                         "field": "bk_inst_id",
                         "operator": "in",
                         "value": [2,4,6]
                    },
                    {
                        "field": "bk_asst_id",
                        "operator": "equal",
                        "value": 3
                    }
                ]
            }
        ]
    },
    "fields":[
        "bk_inst_id",
        "bk_asst_inst_id",
        "bk_asst_obj_id",
        "bk_asst_id",
        "bk_obj_asst_id"
    ],
    "page":{
        "start":0,
        "limit":500
    }
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {
        "info": [
            {
                "bk_inst_id": 2,
                "bk_asst_inst_id": 8,
                "bk_asst_obj_id": "host",
                "bk_asst_id": "connect",
                "bk_obj_asst_id": "bk_switch_connect_host"
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name | Type  | Description                                                         |
|------|-------|---------------------------------------------------------------------|
| info | array | Map array format, returning instance data that meets the conditions |

#### info

| Name            | Type   | Description                                  |
|-----------------|--------|----------------------------------------------|
| bk_inst_id      | int    | Source model instance id                     |
| bk_asst_inst_id | int    | Target model instance id                     |
| bk_asst_obj_id  | string | Associated object model id                   |
| bk_asst_id      | string | Associated type id                           |
| bk_obj_asst_id  | string | Automatically generated model association id |
