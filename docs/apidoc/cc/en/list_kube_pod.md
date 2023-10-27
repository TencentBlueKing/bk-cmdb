### Functional description

list kube pod (version: v3.12.1+, auth: biz access)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                                                                                             |
|-----------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | yes      | biz id                                                                                                                                  |
| filter    | object | no       | pod query filter                                                                                                                        |
| fields    | array  | yes      | pod attribute list, controls which fields in the pod will be returned, can speed up the request and reduce network traffic transmission |
| page      | object | yes      | paging info                                                                                                                             |

#### Filter

This parameter is the filter rule to search for pod based on its attribute fields. This parameter supports the following
two filter rules types. The combined filter rules can be nested with the maximum nesting level of 2. The specific
supported filter rule types are as follows:

##### Combined Filter Rule

This filter rule type defines filter rules composed of other rules, the combined rules support logic and/or
relationships

| Field     | Type   | Required | Description                                                                |
|-----------|--------|----------|----------------------------------------------------------------------------|
| condition | string | yes      | query criteria, support `AND` and `OR`                                     |
| rules     | array  | yes      | query rules, can be of `Combined Filter Rule` or `Atomic Filter Rule` type |

##### Atomic Filter Rule

This filter rule type defines basic filter rules, which represent rules for filtering a field. Any filter rule is either
directly an atomic filter rule, or a combination of multiple atomic filter rules

| Field    | Type                                                                 | Required | Description                                                                                                          |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | yes      | pod's field                                                                                                          |
| operator | string                                                               | yes      | operator, optional values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | different fields and operators correspond to different value formats | yes      | operand                                                                                                              |

Assembly rules can refer to: <https://github.com/Tencent/bk-cmdb/blob/master/src/pkg/filter/README.md>

#### Page

| Field        | Type   | Required | Description                                                                                                                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | yes      | Record start position                                                                                                                                                                                              |
| limit        | int    | yes      | Limit per page, maximum 500                                                                                                                                                                                        |
| sort         | string | no       | Sort the field                                                                                                                                                                                                     |
| enable_count | bool   | yes      | The flag defining Whether to get the the number of query objects. If this flag is true, then the request is to get the quantity. The remaining fields must be initialized, start is 0, and limit is: 0, sort is "" |

**Note:**

- `enable_count`If this flag is true, this request is a get quantity. The remaining fields must be initialized, start is
  0, and limit is: 0, sort is "."
- Paging parameters must be set, and the maximum query data at one time does not exceed 500.

### Request Parameters Example

#### Query Detail Request Parameters Example

```json
{
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
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

#### Query Quantity Request Parameters Example

```json
{
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
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

### Return Result Example

#### Query Detail Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
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

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | object | Data returned by request                                                                |

#### data

| Field | Type  | Description                                               |
|-------|-------|-----------------------------------------------------------|
| count | int   | Number of pods                                            |
| info  | array | Pod list, only returns the fields that is set in `fields` |

#### info[x]

| Field          | Type         | Description                                                                                                                                                                                         |
|----------------|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| name           | string       | pod name                                                                                                                                                                                            |
| priority       | int          | pod priority                                                                                                                                                                                        |
| labels         | string map   | pod labels, key and value are all string, official documentation: http://kubernetes.io/docs/user-guide/labels                                                                                       |
| ip             | string       | pod ip                                                                                                                                                                                              |
| ips            | object array | pod ip list, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#podip-v1-core                                                                                             |
| volumes        | object array | pod volume info list, official documentation: https://kubernetes.io/zh/docs/concepts/storage/volumes/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volume-v1-core |
| qos_class      | enum         | quality of service class, official documentation: https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/quality-service-pod/                                                               |
| node_selectors | string map   | node label selectors, key and value are all string, official documentation: https://kubernetes.io/zh/docs/concepts/scheduling-eviction/assign-pod-node/                                             |
| tolerations    | object array | pod toleration list, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#toleration-v1-core                                                                                |
| operator       | string array | pod operator                                                                                                                                                                                        |
| containers     | object array | container information list                                                                                                                                                                          |
