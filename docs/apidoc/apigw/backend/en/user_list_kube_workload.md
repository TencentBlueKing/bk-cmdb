### Description

Query Workload (Version: v3.12.1+, Permission: Business access)

### Parameters

| Name      | Type   | Required | Description                                                                                                                                                                                     |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                                                                                                     |
| kind      | string | Yes      | Workload type, currently supported workload types are deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods (created directly without going through workload) |
| filter    | object | No       | Workload query conditions                                                                                                                                                                       |
| fields    | array  | No       | Attribute list, control which fields are present in the returned result to speed up API requests and reduce network traffic                                                                     |
| page      | object | Yes      | Pagination information                                                                                                                                                                          |

#### filter Field Description

Filter rules for workload properties, used to search data based on workload properties. This parameter supports two
types of filter rule types, and combination filter rules can be nested, with a maximum nesting of 2 levels. The specific
supported filter rule types are as follows:

##### Combination Filter Rule

A filter rule composed of other rules, supporting logical AND/OR relationships between combined rules.

| Name      | Type   | Required | Description                                                                        |
|-----------|--------|----------|------------------------------------------------------------------------------------|
| condition | string | Yes      | Combination query conditions, supports `AND` and `OR`                              |
| rules     | array  | Yes      | Query rules, can be either `Combination Filter Rule` or `Atomic Filter Rule` types |

##### Atomic Filter Rule

Basic filter rule, represents a rule to filter a specific field. Any filter rule is directly an atomic filter rule or is
composed of multiple atomic filter rules.

| Name     | Type                                                                 | Required | Description                                                                                                                      |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | Yes      | Field of the container                                                                                                           |
| operator | string                                                               | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | Different fields and operators correspond to different value formats | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can be referred
to: [Filter README](https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md)

#### page Field Description

| Name         | Type   | Required | Description                                                                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | Yes      | Record start position                                                                                                                                              |
| limit        | int    | Yes      | Number of records per page, maximum 500                                                                                                                            |
| sort         | string | No       | Sorting field                                                                                                                                                      |
| enable_count | bool   | Yes      | Flag for whether this request is for obtaining the quantity or details. If this flag is true, other fields must be initialized. start is 0, limit is 0, sort is "" |

**Note:**

- `enable_count` If this flag is true, it means this request is to obtain the quantity. At this time, other fields must
  be initialized, start is 0, limit is 0, sort is "".
- Pagination parameters must be set, and the maximum number of queried data at a time should not exceed 500.

### Request Example

#### Detailed Information Request Parameter Example

```json
{
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

#### Quantity Request Parameter Example

```json
{
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

### Response Example

#### Detailed Information Response Example

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
}
```

#### Quantity Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 100,
    "info": []
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
| data       | object | Data returned by the request                                       |

#### data

| Name  | Type  | Description                           |
|-------|-------|---------------------------------------|
| count | int   | Number of records                     |
| info  | array | Actual data, details of the workloads |

#### info[x]

| Name                    | Type   | Description                                                                                                                                     |
|-------------------------|--------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| name                    | string | Workload name                                                                                                                                   |
| labels                  | map    | Labels                                                                                                                                          |
| selector                | object | Workload selector                                                                                                                               |
| replicas                | int    | Number of workload instances                                                                                                                    |
| strategy_type           | string | Workload update mechanism                                                                                                                       |
| min_ready_seconds       | int    | Minimum readiness time for newly created Pods in the absence of any container crashes, only Pods that exceed this time are considered available |
| rolling_update_strategy | object | Rolling update strategy                                                                                                                         |

#### selector

| Name              | Type  | Description           |
|-------------------|-------|-----------------------|
| match_labels      | map   | Match based on labels |
| match_expressions | array | Matching expressions  |

#### match_expressions[x]

| Name     | Type   | Description                                                                                                                   |
|----------|--------|-------------------------------------------------------------------------------------------------------------------------------|
| key      | string | Key of the label                                                                                                              |
| operator | string | Operator, optional values are "In", "NotIn", "Exists", "DoesNotExist"                                                         |
| values   | array  | String array, if the operator is "In" or "NotIn", it must not be empty. If it is "Exists" or "DoesNotExist", it must be empty |

#### rolling_update_strategy

When strategy_type is RollingUpdate, it is not empty; otherwise, it is empty.

| Name            | Type   | Description         |
|-----------------|--------|---------------------|
| max_unavailable | object | Maximum unavailable |
| max_surge       | object | Maximum surge       |

#### max_unavailable

| Name    | Type   | Description                                                                 |
|---------|--------|-----------------------------------------------------------------------------|
| type    | int    | Optional values are 0 (represents int type) or 1 (represents string type)   |
| int_val | int    | When type is 0 (represents int type), it corresponds to the int value       |
| str_val | string | When type is 1 (represents string type), it corresponds to the string value |

#### max_surge

| Name    | Type   | Description                                                                 |
|---------|--------|-----------------------------------------------------------------------------|
| type    | int    | Optional values are 0 (represents int type) or 1 (represents string type)   |
| int_val | int    | When type is 0 (represents int type), it corresponds to the int value       |
| str_val | string | When type is 1 (represents string type), it corresponds to the string value |

**Note:**

- If this request is to query detailed information, count is 0; if it is to query the quantity, info is empty.
