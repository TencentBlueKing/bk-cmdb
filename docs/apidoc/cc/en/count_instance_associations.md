### Functional description

Model instance relation Qty query (v3.10.1+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|    Field    | Type | Required | Description                                                                                                 |
|------------|---------|------|-----------------------------------------------------------------------------------------------------------------|
| bk_biz_id  | int |no| Business ID, which needs to be provided when querying mainline model                                                                              |
| bk_obj_id  | string  |yes| Model ID                                                                                                          |
| conditions | object  |no| Combined query criteria: AND and OR are supported for combination, and can be nested up to 3 layers. Each layer supports 20 OR criteria at most. If this parameter is not specified, it means all matches (i.e., Contexts are null).|

#### conditions

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition |  string  |yes    | Rule operator|
| rules |  array  |yes     | Scope condition rule for selected business|

#### conditions.rules

|   Field   | Type| Required| Description                                                                                                     |
|----------|--------|------|-----------------------------------------------------------------------------------------------------------|
| field    |  string | yes      | Condition field, optional value id, bk_inst_id, bk_obj_id, bk_Asst_inst_id, bk_Asst_obj_id, bk_obj_Asst_id, bk_Asst_id   |
| operator | string |yes| Operator, optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between, etc|
| value    |   -    |no| The expected value of the condition field. Different values correspond to different value formats. The array type value supports up to 500 elements                          |

For details of assembly rules, please refer to: https: //github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

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
        "count": 1
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

| Field|   Type| Description                       |
|-------|---------|----------------------------|
| count | int |Returns the number of instance data that meets the condition|
