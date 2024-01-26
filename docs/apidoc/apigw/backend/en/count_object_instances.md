### Description

Universal model instance quantity query (v3.10.1+)

### Parameters

| Name           | Type   | Required | Description                                                                                                                                                                                               |
|----------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id      | string | Yes      | Model ID                                                                                                                                                                                                  |
| conditions     | object | No       | Combined query conditions, supporting AND and OR, can be nested, up to 3 layers of nesting, up to 20 OR conditions per layer, not specifying this parameter means matching all (i.e., conditions is null) |
| time_condition | object | No       | Query conditions for model instances by time                                                                                                                                                              |

#### conditions

| Name      | Type   | Required | Description                                     |
|-----------|--------|----------|-------------------------------------------------|
| condition | string | Yes      | Rule operator                                   |
| rules     | array  | Yes      | Range conditions for the selected business rule |

#### conditions.rules

| Name     | Type   | Required | Description                                                                                                                                                     |
|----------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Condition field                                                                                                                                                 |
| operator | string | Yes      | Operator, optional values equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between, etc.                              |
| value    | -      | No       | Expected value of the condition field. Different operators correspond to different value formats. The maximum number of elements for an array-type value is 500 |

For detailed assembly rules, please refer
to: [QueryBuilder](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### time_condition

| Name  | Type   | Required | Description                              |
|-------|--------|----------|------------------------------------------|
| oper  | string | Yes      | Operator. Currently only supports 'and'. |
| rules | array  | Yes      | Time query conditions                    |

#### time_condition.rules

| Name  | Type   | Required | Description                                              |
|-------|--------|----------|----------------------------------------------------------|
| field | string | Yes      | Field name, which should be the field name of the model. |
| start | string | Yes      | Start time in the format yyyy-MM-dd hh:mm:ss             |
| end   | string | Yes      | End time in the format yyyy-MM-dd hh:mm:ss               |

### Request Example

```json
{
    "bk_obj_id":"bk_switch",
    "conditions":{
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
    }
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "count": 1
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |

#### data

| Name  | Type | Description                                  |
|-------|------|----------------------------------------------|
| count | int  | Number of instances that meet the conditions |
