### Function Description

Query workload (version: v3.10.23+, auth: none)

### Request parameters

{{ common_args_desc }}

#### Interface parameters
- Common fields.

| field | type | mandatory | description |
|----------------------------|------------|--------|--------------------------------------------|
|bk_biz_id | int| yes |business id|
|kind | string | yes |workload type, the current built-in workload types are deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods (put those who do not pass workload but directly create Pod)|
| filter | object | no | query condition |
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

### Example request parameters
#### Query Detail Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "kind": "deployment",
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "cluster_uid",
                "operator": "equal",
                "value": "1"
            },
            {
                "field": "namespace",
                "operator": "equal",
                "value": "namespace1"
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

#### Query Quantity Request Parameters Example
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "kind": "deployment",
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "cluster_uid",
                "operator": "equal",
                "value": "1"
            },
            {
                "field": "namespace",
                "operator": "equal",
                "value": "namespace1"
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
        "enable_count": true
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
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
    "data": {
        "count": 10,
        "info": []
    }
}
```

### Return result parameter description
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

| Field | Type | Description |
| ----- | ----- | ------------ |
|name | string |workload name|
| labels| map |label|
| selector| object | workload selector|
| replicas| no |number of workload instances|
| strategy_type| string |workload update mechanism|
| min_ready_seconds| int |Specifies the minimum time a newly created Pod will be ready without any container crashes, and only after that time will the Pod be considered available|
| rolling_update_strategy| object | rolling update strategy|

#### selector
| field | type | description |
| ----- | ----- | ------------ |
|match_labels | map |match_by_label|
|match_expressions | array |match expressions |

#### match_expressions[0]
| field | type | description |
| ----- | ----- | ------------ |
|key | string |key of tag|
|operator | string |operator, optional values: "In", "NotIn", "Exists", "DoesNotExist"|
|values | array | array of strings, cannot be empty if the operator is "In" or "NotIn", must be empty if it is "Exists" or "DoesNotExist"|

#### rolling_update_strategy
When strategy_type is RollingUpdate, it is not empty, otherwise it is empty.

| field | type | description |
| ----- | ----- | ------------ |
|max_unavailable | object |max_unavailable | object
|max_surge | object |max_overflow|

#### max_unavailable
| field | type | description |
| ----- | ----- | ------------ |
|type | int |Optional values of 0 (for int types) or 1 (for string types)|
|int_val | int |When type is 0 (for int type), the corresponding int value |
|str_val | string |when type is 1(string type),the corresponding string value|

#### max_surge
| field | type | description |
| ----- | ----- | ------------ |
|type | int | optional value of 0 (for int type) or 1 (for string type) |
|int_val | int |When type is 0 (for int type), the corresponding int value |
|str_val | string |when type is 1(string type),the corresponding string value|

**Note:**
- If this request is to query details, count is 0. If the query is quantity, info is empty.
