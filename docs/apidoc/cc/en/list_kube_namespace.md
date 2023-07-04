### Functional description

Query namespace(version: v3.10.23+, auth: none)

### Request parameters

{{ common_args_desc }}

#### Interface parameters
- common fields.

| field | type | required | description |
|----------------------------|------------|--------|--------------------------------------------|
| bk_biz_id | int| Yes | business id|
| filter | object | no | query criteria |
| fields | array | No | A list of attributes that control which fields are returned in the result, to speed up interface requests and reduce network traffic transfer |
| page | object | yes | paging information |

#### filter
This parameter is the filter rule to search for container based on its attribute fields. This parameter supports the following two filter rules types. The combined filter rules can be nested with the maximum nesting level of 2. The specific supported filter rule types are as follows:

##### combined filter rule

This filter rule type defines filter rules composed of other rules, the combined rules support logic and/or relationships

| Field     | Type   | Required | Description             |
|----------|-------|-----|---------------------- -------------|
| condition | string | yes | query criteria, support `AND` and `OR` |
| rules | array | yes |  query rules, can be of `combined filter rule` or `atomic filter rule` type |

##### atomic filter rule

This filter rule type defines basic filter rules, which represent rules for filtering a field. Any filter rule is either directly an atomic filter rule, or a combination of multiple atomic filter rules

| Field    | Type                                                                 | Required | Description                                                                                                          |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | yes      | container's field                                                                                                    |
| operator | string                                                               | yes      | operator, optional values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | different fields and operators correspond to different value formats | yes      | operand                                                                                                              |

Assembly rules can refer to: <https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md>

#### page

| field | type | required | description |
| ----- | ------ | ---- | -------------------- |
| start | int | Yes | Record start position |
| limit | int | Yes | Limit the number of entries per page, up to 500 |
| sort | string | No | Sort field |
| enable_count | bool | Yes | A flag for whether to get the number of query objects. If this flag is true then the request is to get the number, the rest of the fields must be initialized, start is 0, limit is :0, sort is "" |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_cluster_id",
                "operator": "equal",
                "value": 1
            },
            {
                "field": "name",
                "operator": "equal",
                "value": "test"
            }
        ]
    },
    "fields": [
        "name"
    ],
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "name",
        "enable_count": false
    }
}
```

- The second, using a unique identifier in cc.

| field | type | required | description |
| ----------------------------|------------|--------|--------------------------------------------|
| bk_cluster_id | int| No |cluster's unique id in cc|

### Request Parameters Example
#### Query Detail Request Parameters Example
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "bk_cluster_id": 1,
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "name",
                "operator": "equal",
                "value": "test"
            }
        ]
    },
    "fields": [
        "name"
    ],
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "name",
        "enable_count": false
    }
}
```
#### Query Quantity Request Parameters Example
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_cluster_id",
                "operator": "equal",
                "value": 1
            },
            {
                "field": "name",
                "operator": "equal",
                "value": "test"
            }
        ]
    },
    "page": {
        "enable_count":true
    }
}
```

### Return Result Example
#### Query Detail Return Result Example
```json

{
    "result": true,
    "code": 0,
    "data": {
        "count": 0,
        "info": [
            {	
                "name": "test"
            }
        ]
    },
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```
#### Query Quantity Return Result Example
```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":100,
        "info":[
        ]
    },
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return Result Parameters Description

#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | Whether the request was successful or not. true:request successful; false request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | data returned by the request

#### data
| field | type | description |
| ----- | ----- | ------------ |
| count | int | Number of records |
| info | array | The actual data, returning only the fields set in fields |

#### data.info
| field | type | description |
| ----- | ----- | ------------ |
| name | string | namespace name|
| labels| map | labels|
| resource_quotas| array | namespace CPU and memory requests and limits|

#### resource_quotas[0]
| field | type | description |
| ----- | ----- | ------------ |
| hard | object | hard limits required per named resource|
|scopes | array |Quota scopes, with optional values of "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", "CrossNamespacePodAffinity"|.
|scope_selector | object |scope selector|

#### scope_selector
| field | type | description |
| ----- | ----- | ------------ |
|match_expressions | array | match_expressions |

#### match_expressions[0]
| field | type | description |
| ----- | ----- | ------------ |
|scope_name | array |Quota scope, optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBe
|operator | string  |selector operator，optional values are："In"、"NotIn"、"Exists"、"DoesNotExist"|
|values | array |string array，if the operator is "In"or "NotIn",can not be null，if the operator is "Exists" or "DoesNotExist"，it must be null|

**Note:**
- If this request is to query details, count is 0. If the query is quantity, info is empty.
