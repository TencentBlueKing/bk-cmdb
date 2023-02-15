### Functional description

list kube container (version: v3.10.23+, auth: None)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field           | Type   | Required | Description                                                                                                                                         |
|-----------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id       | int    | yes      | biz id                                                                                                                                              |
| filter          | object | no       | container query filter                                                                                                                              |
| fields          | array  | yes      | container attribute list, controls which fields in the container will be returned, can speed up the request and reduce network traffic transmission |
| page            | object | yes      | paging info                                                                                                                                         |

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

| Field        | Type   | Required | Description                                                                                                                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | yes      | Record start position                                                                                                                                                                                              |
| limit        | int    | yes      | Limit per page, maximum 500                                                                                                                                                                                        |
| sort         | string | no       | Sort the field                                                                                                                                                                                                     |
| enable_count | bool   | yes      | The flag defining Whether to get the the number of query objects. If this flag is true, then the request is to get the quantity. The remaining fields must be initialized, start is 0, and limit is: 0, sort is "" |

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
                "value": "container1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "uid",
                        "operator": "not_in",
                        "value": [
                            "xxxxxx"
                        ]
                    },
                    {
                        "field": "image",
                        "operator": "equal",
                        "value": "ccr.ccs.tencentyun.com/library/coredns:1.6.2"
                    }
                ]
            }
        ]
    },
    "fields": [
        "name",
        "uid"
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
                "value": "container1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "uid",
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
                "name": "container2",
                "uid": "xxx"
            },
            {
                "name": "container3",
                "uid": "xxx"
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

| Field | Type  | Description                                                     |
|-------|-------|-----------------------------------------------------------------|
| count | int   | Number of containers                                            |
| info  | array | Container list, only returns the fields that is set in `fields` |

#### info
| Field       | Type         | Description                                                                                                                                                                                                                                     |
|-------------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| name        | string       | container name                                                                                                                                                                                                                                  |
| uid         | string       | container uid                                                                                                                                                                                                                                   |
| image       | string       | container image                                                                                                                                                                                                                                 |
| ports       | object array | container port information list, format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core                                                                                                             |
| args        | string array | start arguments                                                                                                                                                                                                                                 |
| started     | timestamp    | start time                                                                                                                                                                                                                                      |
| limits      | object       | resource limits, official documentation: https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/                                                                                                                                         |
| requests    | object       | resource requests, official documentation: https://kubernetes.io/zh/docs/concepts/policy/resource-quotas/                                                                                                                                       |
| liveness    | object       | liveness probe, official documentation: https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core   |
| environment | object array | environment variables, official documentation: https://kubernetes.io/zh/docs/tasks/inject-data-application/define-environment-variable-container/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envvar-v1-core |
| mounts      | object array | volume mounts, official documentation: https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-volume-storage/ , format: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumemount-v1-core                 |
