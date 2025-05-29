### Description

Query namespace (Version: v3.12.1+, Permission: Business access)

### Parameters

| Name      | Type   | Required | Description                                                                                                               |
|-----------|--------|----------|---------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                               |
| filter    | object | No       | Namespace query conditions                                                                                                |
| fields    | array  | No       | Property list, controls which fields are returned in the result to speed up interface requests and reduce network traffic |
| page      | object | Yes      | Pagination information                                                                                                    |

#### filter Field Description

Filter rules for namespace properties, used to search for data based on namespace properties. This parameter supports
two types of filter rule types, where combination filter rules can be nested, and at most 2 levels of nesting. The
specific supported filter rule types are as follows:

##### Combination filter rules

Filter rules composed of other rules, supporting logical AND/OR relationships between rules

| Name      | Type   | Required | Description                                                                     |
|-----------|--------|----------|---------------------------------------------------------------------------------|
| condition | string | Yes      | Combined query condition, supports both `AND` and `OR`                          |
| rules     | array  | Yes      | Query rules, can be of type `Combination filter rules` or `Atomic filter rules` |

##### Atomic filter rules

Basic filter rules, indicating the rules for filtering a field. Any filter rule is directly an atomic filter rule or is
composed of multiple atomic filter rules.

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Namespace field                                                                                                                  |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can be referred
to: [Filter README](https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md)

#### page

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

#### Detailed Information Request Parameter

```json
{
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

#### Quantity Request Example

```json
{
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
    "enable_count": true
  }
}
```

### Response Example

#### Detailed Information Interface Response

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
    "info": [
    ]
  },
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

| Name  | Type  | Description                                            |
|-------|-------|--------------------------------------------------------|
| count | int   | Number of records                                      |
| info  | array | Actual data, only returns the fields set in the fields |

#### info[x]

| Name            | Type   | Description                                  |
|-----------------|--------|----------------------------------------------|
| name            | string | Namespace name                               |
| labels          | map    | Labels                                       |
| resource_quotas | array  | Namespace CPU and memory requests and limits |

#### resource_quotas[x]

| Name           | Type   | Description                                                                                                                                    |
|----------------|--------|------------------------------------------------------------------------------------------------------------------------------------------------|
| hard           | object | Hard limit for each named resource                                                                                                             |
| scopes         | array  | Quota scope, optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", "CrossNamespacePodAffinity" |
| scope_selector | object | Scope selector                                                                                                                                 |

#### scope_selector

| Name              | Type  | Description       |
|-------------------|-------|-------------------|
| match_expressions | array | Match expressions |

#### match_expressions[x]

| Name       | Type   | Description                                                                                                                                    |
|------------|--------|------------------------------------------------------------------------------------------------------------------------------------------------|
| scope_name | array  | Quota scope, optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", "CrossNamespacePodAffinity" |
| operator   | string | Selector operator, optional values are: "In", "NotIn", "Exists", "DoesNotExist"                                                                |
| values     | array  | String array. If the operator is "In" or "NotIn", it cannot be empty. If it is "Exists" or "DoesNotExist", it must be empty.                   |
