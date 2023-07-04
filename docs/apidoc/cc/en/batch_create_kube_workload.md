### Function description

Batch create workload (version: v3.10.23+, auth: create container workload)

### Request parameters

{{ common_args_desc }}

#### Interface parameters

| field | type | required | description |
|----------------------------|------------|--------|--------------------------------------------|
|bk_biz_id | int| yes |business id|
|kind | string | yes |workload type, the current built-in workload types are deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods (put those that do not pass the workload but directly create Pod)|
| data | array| Yes | array, limit to 200 at a time|

#### data

| field | type | required | description |
|----------------------------|------------|--------|--------------------------------------------|
|bk_namespace_id | int |yes |namespace's unique identifier in cc|
|name | string |yes |workload name|
|labels| map |no |labels|
| selector| object | no |workload selector|
| replicas| no | no |number of workload instances|
| strategy_type| string | no |workload update mechanism|
| min_ready_seconds| int | No | Specifies the minimum time that a newly created Pod will be ready without any container crashes, and only after that time will the Pod be considered available|
| rolling_update_strategy| object | No | Rolling update strategy|

#### selector
| field | type | required | description |
| ----- | ----- | ------------|------------ |
|match_labels | map | no| match by label|
|match_expressions | array |no|match_expressions|

#### match_expressions[0]
| field | type | required | description |
| ----- | ----- | ------------|------------ |
|key | string | is the key| of the |tag
|operator | string | is the |operator, with optional values: "In", "NotIn", "Exists", "DoesNotExist"|
|values | array |no| Array of strings, cannot be empty if the operator is "In" or "NotIn", must be empty if it is "Exists" or "DoesNotExist"|

#### rolling_update_strategy
When strategy_type is RollingUpdate, it is not empty, otherwise it is empty.

| field | type | mandatory | description |
| ----- | ----- | ------------|------------ |
|max_unavailable | object |no|max_unavailable|
|max_surge | object |no|max_overflow|

#### max_unavailable
| field | type | mandatory | description |
| ----- | ----- | ------------|------------ |
|type | int |Yes|Optional value of 0 (for int type) or 1 (for string type)|
|int_val | int |No|When type is 0 (for int type), it cannot be null, and the corresponding int value|
|str_val | string |no|when type is 1(for string type),cannot be null,corresponding string value|

#### max_surge
| field | type | mandatory | description |
| ----- | ----- | ------------|------------ |
|type | int | yes | optional value of 0 (for int type) or 1 (for string type) |
|int_val | int |No|When type is 0 (for int type), it cannot be null, and the corresponding int value|
|str_val | string |no|When type is 1 (for string type), it cannot be empty, and the corresponding string value|

### Example request parameters
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "kind": "deployment",
    "data": [
      {
        "bk_namespace_id": 1,
        "name": "test",
        "labels": {
            "test": "test",
            "test2": "test2"  
        },
        "selector": {
            "match_labels": {
                "test": "test",
                "test2": "test2" 
            },
            "match_expressions": [
                {
                    "key": "tier",
                    "operator": "In", 
                    "values": ["cache"]
                }
            ]
        },
        "replicas": 1,
        "strategy_type": "RollingUpdate",
        "min_ready_seconds": 1,
        "rolling_update_strategy": {
            "max_unavailable": {
                "type": 0,
                "int_val": 1
            },
            "max_surge": {
                "type": 0,
                "int_val": 1
            }
        }
      }  
    ]   
}
```

### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "data": {
      "ids": [1]
    },
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```
**Note:**
- The order of the node ID array in the returned data is consistent with the order of the array data in the parameter.

### Return result parameter description
#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | Whether the request was successful or not. true:request successful; false request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | The data returned by the request.|

#### data

| field | type | description |
|----------- |-----------|----------|
| ids | array | Array of unique identifiers in cc |
