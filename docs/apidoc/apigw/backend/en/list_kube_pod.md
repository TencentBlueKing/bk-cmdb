### Description

Query Pod list (Version: v3.12.1+, Permission: Business access)

### Parameters

| Name      | Type   | Required | Description                                                                                                                         |
|-----------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                                         |
| filter    | object | No       | Conditions for querying pods                                                                                                        |
| fields    | array  | Yes      | List of pod properties, control which fields are present in the returned result to speed up API requests and reduce network traffic |
| page      | object | Yes      | Pagination information                                                                                                              |

#### filter Field Description

Filter rules for pod properties, used to search data based on pod properties. This parameter supports two types of
filter rule types, and combination filter rules can be nested, with a maximum nesting of 2 levels. The specific
supported filter rule types are as follows:

##### Combination Filter Rule

A filter rule composed of other rules, supporting logical AND/OR relationships between combined rules.

| Name      | Type   | Required | Description                                                                                  |
|-----------|--------|----------|----------------------------------------------------------------------------------------------|
| condition | string | Yes      | Combination query conditions, supports `AND` and `OR`                                        |
| rules     | array  | Yes      | Query rules, can be either `Combination Query Parameters` or `Atomic Query Parameters` types |

##### Atomic Filter Rule

Basic filter rule, represents a rule to filter a specific field. Any filter rule is directly an atomic filter rule or is
composed of multiple atomic filter rules.

| Name     | Type                                                                 | Required | Description                                                                                                                      |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | Yes      | Field of the pod                                                                                                                 |
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
  "bk_biz_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "pod1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "priority",
            "operator": "not_in",
            "value": [
              2,
              6
            ]
          },
          {
            "field": "qos_class",
            "operator": "equal",
            "value": "Burstable"
          }
        ]
      }
    ]
  },
  "fields": [
    "name",
    "priority"
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
  "bk_biz_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "pod1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "priority",
            "operator": "not_in",
            "value": [
              2,
              6
            ]
          },
          {
            "field": "qos_class",
            "operator": "equal",
            "value": "Burstable"
          }
        ]
      }
    ]
  },
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
  "message": "success",
  "permission": null,
  "data": {
    "count": 0,
    "info": [
      {
        "name": "pod2",
        "priority": 1
      },
      {
        "name": "pod3",
        "priority": 5
      }
    ]
  }
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
    "count": 10,
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

| Name  | Type  | Description                      |
|-------|-------|----------------------------------|
| count | int   | Number of records                |
| info  | array | Actual data, details of the pods |

#### info[x]

| Name           | Type         | Description                    |
|----------------|--------------|--------------------------------|
| name           | string       | Name                           |
| priority       | int          | Priority                       |
| labels         | string map   | Labels                         |
| ip             | string       | Container network IP           |
| ips            | object array | Array of container network IPs |
| volumes        | object array | Used volume information        |
| qos_class      | enum         | Quality of service             |
| node_selectors | string map   | Node label selector            |
| tolerations    | object array | Tolerations                    |
| operator       | string array | Pod operator                   |
| containers     | object array | Container data                 |
