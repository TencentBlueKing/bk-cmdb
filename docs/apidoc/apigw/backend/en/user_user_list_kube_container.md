### Description

Query Container list (Version: v3.12.1+, Permission: Business access)

### Parameters

| Name      | Type         | Required | Description                                                                                                                                    |
|-----------|--------------|----------|------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int          | Yes      | Business ID                                                                                                                                    |
| bk_pod_id | int          | No       | ID of the pod to which the container belongs                                                                                                   |
| filter    | object       | No       | Query conditions for the container                                                                                                             |
| fields    | string array | Yes      | List of container properties, controls which fields are returned in the container, speeding up interface requests and reducing network traffic |
| page      | object       | Yes      | Pagination information                                                                                                                         |

#### filter Field Description

Filter rules for container properties, used to search for data based on container properties. This parameter supports
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
| field    | string | Yes      | Container field                                                                                                                  |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

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
  "bk_pod_id": 4,
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "container1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "container_uid",
            "operator": "not_in",
            "value": [
              "xxxxxx"
            ]
          },
          {
            "field": "image",
            "operator": "equal",
            "value": "xxx"
          }
        ]
      }
    ]
  },
  "fields": [
    "name",
    "container_uid"
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
        "value": "container1"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "container_uid",
            "operator": "not_in",
            "value": [
              "xxxxxx"
            ]
          },
          {
            "field": "image",
            "operator": "equal",
            "value": "xxx"
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
        "name": "container2",
        "container_uid": "xxx"
      },
      {
        "name": "container3",
        "container_uid": "xxx"
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

| Name  | Type  | Description           |
|-------|-------|-----------------------|
| count | int   | Number of records     |
| info  | array | Actual container data |

#### info[x]

| Name          | Type         | Description                                                                                                                                                                                                                                                                                                 |
|---------------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| name          | string       | Name                                                                                                                                                                                                                                                                                                        |
| container_uid | string       | Container UID                                                                                                                                                                                                                                                                                               |
| image         | string       | Image information                                                                                                                                                                                                                                                                                           |
| ports         | object array | Container ports, format: [ContainerPort](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core)                                                                                                                                                                        |
| args          | string array | Startup parameters                                                                                                                                                                                                                                                                                          |
| started       | timestamp    | Startup time                                                                                                                                                                                                                                                                                                |
| limits        | object       | Resource limits, official documentation: Resource Quotas                                                                                                                                                                                                                                                    |
| requests      | object       | Resource request size, official documentation: Resource Quotas                                                                                                                                                                                                                                              |
| liveness      | object       | Liveness probe, official documentation: [Configure Liveness, Readiness, and Startup Probes](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) , format: [Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core) |
| environment   | object array | Environment variables, official documentation: [Define Environment Variable Container](https://kubernetes.io/zh/docs/tasks/inject-data-application/define-environment-variable-container/) , format: [EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envvar-v1-core)          |
| mounts        | object array | Volume mounts, official documentation: [Configure Volume Storage](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-volume-storage/) , format: [VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumemount-v1-core)                                  |
