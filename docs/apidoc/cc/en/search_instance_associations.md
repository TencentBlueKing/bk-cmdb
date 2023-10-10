### Functional description

General model instance relation query (v3.10.1+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|    Field    | Type | Required | Description                                                                                                 |
|------------|---------|------|-----------------------------------------------------------------------------------------------------------------|
| bk_biz_id  | int |no| Business ID, which is required for mainline model query                                                                              |
| bk_obj_id  | string  |yes| Model ID                                                                                                          |
| conditions | object  |no| Combined query criteria: AND and OR are supported for combination, and can be nested up to 3 layers. Each layer supports 20 OR criteria at most. If this parameter is not specified, it means all matches (i.e., Contexts are null).|
| fields     |  array   | no | Specify the fields to be returned. Fields that do not exist will be ignored. If not specified, all fields will be returned (returning all fields will affect performance, and it is recommended to return on demand).    |
| page       |  object  |yes| Paging settings                                                                                                        |

#### conditions

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  |yes     | Scope condition rule for selected business|

#### conditions.rules

|   Field   | Type| Required| Description                                                                                                     |
|----------|--------|------|-----------------------------------------------------------------------------------------------------------|
| field    |  string |yes| Condition field, optional value id, bk_inst_id, bk_obj_id, bk_Asst_inst_id, bk_Asst_obj_id, bk_obj_Asst_id, bk_Asst_id   |
| operator | string |yes| Operator, optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between, etc|
| value    |   -    |no| The expected value of the condition field. Different values correspond to different value formats. The array type value supports up to 500 elements                          |

For details of assembly rules, please refer to https: //github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

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

### Return result parameter

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

| Field| Type| Description                                |
|------|-------|-------------------------------------|
| info | array |map array format, which returns instance data that meets the condition|

#### info

| Field| Type| Description                                |
|------|-------|-------------------------------------|
| bk_inst_id | int |Source model instance id|
| bk_asst_inst_id|  int| Target model instance id|
| bk_asst_obj_id|  string| Association relation target model id|
| bk_asst_id|  string| Association type id|
| bk_obj_asst_id|  string| Auto-generated model association id|


