### Description

General model instance query (Version: v3.10.1+, Permission: Model instance query permission)

### Parameters

| Name           | Type   | Required | Description                                                                                                                                                                                                           |
|----------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id      | string | Yes      | Model ID                                                                                                                                                                                                              |
| conditions     | object | No       | Combined query conditions, supports AND and OR, can be nested, up to 3 levels of nesting, each OR condition supports up to 20 conditions, not specifying this parameter means matching all (i.e., conditions is null) |
| time_condition | object | No       | Query conditions for model instances based on time                                                                                                                                                                    |
| fields         | array  | No       | Specify the fields to be returned. Fields that are not available will be ignored. If not specified, all fields will be returned (returning all fields will impact performance, it is recommended to return as needed) |
| page           | object | Yes      | Pagination settings                                                                                                                                                                                                   |

#### conditions

| Name      | Type   | Required | Description                                      |
|-----------|--------|----------|--------------------------------------------------|
| condition | string | Yes      | Rule operator                                    |
| rules     | array  | Yes      | Range conditions for the selected business scope |

#### conditions.rules

| Name     | Type   | Required | Description                                                                                                                                                   |
|----------|--------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Condition field                                                                                                                                               |
| operator | string | Yes      | Operator, optional values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between, etc.                           |
| value    | -      | No       | Expected value of the condition field. Different operators correspond to different value formats. The maximum number of elements for array type values is 500 |

Detailed assembly rules can be found
here: [querybuilder](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### time_condition

| Name  | Type   | Required | Description                             |
|-------|--------|----------|-----------------------------------------|
| oper  | string | Yes      | Operator, currently only supports 'and' |
| rules | array  | Yes      | Time query conditions                   |

#### time_condition.rules

| Name  | Type   | Required | Description                             |
|-------|--------|----------|-----------------------------------------|
| field | string | Yes      | Value is the field name of the model    |
| start | string | Yes      | Start time, format: yyyy-MM-dd hh:mm:ss |
| end   | string | Yes      | End time, format: yyyy-MM-dd hh:mm:ss   |

#### page

| Name  | Type   | Required | Description                                                                                           |
|-------|--------|----------|-------------------------------------------------------------------------------------------------------|
| start | int    | Yes      | Record start position                                                                                 |
| limit | int    | Yes      | Page limit, maximum 500                                                                               |
| sort  | string | No       | Retrieval sorting, follow the MongoDB semantic format {KEY}:{ORDER}, default sorting by creation time |

### Request Example

```json{
    "bk_obj_id": "bk_switch",
    "conditions": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_inst_name",
                "operator": "equal",
                "value": "switch"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                         "field": "bk_inst_id",
                         "operator": "not_in",
                         "value": [2,4,6]
                    },
                    {
                        "field": "bk_inst_id",
                        "operator": "equal",
                        "value": 3
                    }
                ]
            }
        ]
    },
    "time_condition": {
            "oper": "and",
            "rules": [
                {
                    "field": "create_time",
                    "start": "2021-05-13 01:00:00",
                    "end": "2021-05-14 01:00:00"
                }
            ]
    },
    "fields": [
        "bk_inst_id",
        "bk_inst_name"
    ],
    "page": {
        "start": 0,
        "limit": 500
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
                "bk_inst_id": 1,
                "bk_inst_name": "switch-instance"
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| data       | object | Data returned by the request                                      |
| permission | object | Permission information                                            |

#### data

| Name | Type  | Description                                                         |
|------|-------|---------------------------------------------------------------------|
| info | array | Map array format, returning instance data that meets the conditions |

#### info

| Name         | Type   | Description   |
|--------------|--------|---------------|
| bk_inst_id   | int    | Instance ID   |
| bk_inst_name | string | Instance name |
