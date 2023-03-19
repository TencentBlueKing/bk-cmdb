### Functional description

Query model instance according to Association relation instance

- This interface is only applicable to custom hierarchical model and general model instances, and is not applicable to model instances such as business, set, module, host, etc

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
| bk_obj_id           |  string     | yes  | Model ID                      |
| page                |  object     | yes  | Paging parameter                    |
| condition           |  object     | no     | Model instance query criteria with Association relationship                    |
| time_condition      |  object     | no     | Query criteria for querying model instances by time|
| fields              |  object     | no     | Specifies the field returned by the query model instance, where key is the model ID and value is the model attribute field to be returned by the query model|

#### page

| Field      | Type      | Required   | Description                |
|-----------|------------|--------|----------------------|
| start     |   int       | yes  | Record start position         |
| limit     |   int       | yes  | Limit bars per page, Max. 200|
| sort      |   string    | no     | Sort field             |

#### condition
The user in the example is the model

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| field     | string      | yes   | The value is the field name of the model                                               |
| operator  |string      | yes   | Value is: $regex $eq $ne                                           |
| value     | string      | yes   | The value corresponding to the model attribute name of the field configuration                   |

#### time_condition

| Field   | Type   | Required| Description              |
|-------|--------|-----|--------------------|
| oper  | string |yes| Operator, currently only and is supported|
| rules | array  |yes| Time query criteria         |

#### rules

| Field   | Type   | Required| Description                             |
|-------|--------|-----|----------------------------------|
| field | string |yes| The value is the field name of the model                  |
| start | string |yes| Start time in the format yyyy MM dd hh: mm:ss|
| end   |  string |yes| End time in the format yyyy MM dd hh: mm:ss|


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "bk_switch",
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_inst_id"
    },
    "fields": {
        "bk_switch": [
            "bk_asset_id",
            "bk_inst_id",
            "bk_inst_name",
            "bk_obj_id"
        ]
    },
    "condition": {
        "user": [
            {
                "field": "operator",
                "operator": "$regex",
                "value": "admin"
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

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 2,
        "info": [
            {
                "bk_asset_id": "sw00001",
                "bk_inst_id": 1,
                "bk_inst_name": "sw1",
                "bk_obj_id": "bk_switch"
            },
            {
                "bk_asset_id": "sw00002",
                "bk_inst_id": 2,
                "bk_inst_name": "sw2",
                "bk_obj_id": "bk_switch"
            }
        ]
    }
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data

| Field      | Type      | Description         |
|-----------|-----------|--------------|
| count     |  int       | Number of records     |
| info      |  array     | Model instance actual data|
