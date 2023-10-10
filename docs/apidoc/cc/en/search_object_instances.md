### Functional description

General model instance query (v3.10.1+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|    Field    |  Type  | Required | Description                                                                                                 |
|------------|--------|------|-----------------------------------------------------------------------------------------------------------------|
| bk_obj_id  | string |yes| Model ID                                                                                                          |
| conditions | object |no| Combined query criteria. Combination supports AND and OR, and can be nested up to 3 levels. Each level supports up to 20 OR criteria. If this parameter is not specified, it means all matches (i.e. Contexts are null).|
| time_condition      |  object     | no     | Query criteria for querying model instances by time|
| fields     |  array  |no| Specify the fields to be returned. Fields that do not exist will be ignored. If not specified, all fields will be returned (returning all fields will affect performance, and it is recommended to return on demand).    |
| page       |  object |yes| Paging settings                                                                                                        |

#### conditions

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  | yes      | Scope condition rule for selected business|

#### conditions.rules

|   Field   | Type| Required| Description                                                                                                     |
|----------|--------|------|-----------------------------------------------------------------------------------------------------------|
| field    |  string |yes| Condition field                                                                                                  |
| operator | string |yes| Operator, optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between, etc|
| value    |   -    |no| The expected value of the condition field. Different values correspond to different value formats. The array type value supports up to 500 elements                          |

For details of assembly rules, please refer to: https: //github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### time_condition

| Field   | Type   | Required| Description              |
|-------|--------|-----|--------------------|
| oper  | string |yes| Operator, currently only and is supported|
| rules | array  |yes| Time query criteria         |

#### time_condition.rules

| Field   | Type   | Required| Description                             |
|-------|--------|-----|----------------------------------|
| field | string |yes| The value is the field name of the model                  |
| start | string |yes| Start time in the format yyyy MM dd hh: mm:ss|
| end   |  string |yes| End time in the format yyy MM dd hh: mm:ss|

#### page

|  Field| Type| Required| Description                                                            |
|-------|--------|------|------------------------------------------------------------------|
| start | int    | yes | Record start position                                                     |
| limit | int    | yes | Limit bars per page, Max. 500                                            |
| sort  | string |no| Retrieve sort, following mongordb semantic format {KEY}:{ORDER}, sorted by creation time by default|

### Request Parameters Example

```json
{
    "bk_app_code":"code",
    "bk_app_secret":"secret",
    "bk_username": "xxx",
    "bk_token":"xxxx",
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
    },
    "fields":[
        "bk_inst_id",
        "bk_inst_name"
    ],
    "page":{
        "start":0,
        "limit":500
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

### Return result parameter

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field| Type| Description                                |
|------|-------|-------------------------------------|
| info | array |map array format, which returns instance data that meets the condition|

#### info
| Field| Type| Description                                |
|------|-------|-------------------------------------|
| bk_inst_id | int |Instance id|
| bk_inst_name | string |Instance name|
