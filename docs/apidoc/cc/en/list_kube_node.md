### Functional description

list container nodes (v3.12.1+, permission: biz access)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                                                    |
|-----------|--------|----------|------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | yes      | business ID                                                                                    |
| filter    | object | no       | Container node query scope                                                                     |
| fields    | array  | no       | The attribute of the container node to be queried, if not written, it means to search all data |
| page      | object | yes      | Paging condition                                                                               |

#### filter

This parameter is the filter rule to search for container based on its attribute fields. This parameter supports the
following two filter rules types. The combined filter rules can be nested with the maximum nesting level of 2. The
specific supported filter rule types are as follows:

##### combined filter rule

This filter rule type defines filter rules composed of other rules, the combined rules support logic and/or
relationships

| Field     | Type   | Required | Description                                                                |
|-----------|--------|----------|----------------------------------------------------------------------------|
| condition | string | yes      | query criteria, support `AND` and `OR`                                     |
| rules     | array  | yes      | query rules, can be of `combined filter rule` or `atomic filter rule` type |

##### atomic filter rule

This filter rule type defines basic filter rules, which represent rules for filtering a field. Any filter rule is either
directly an atomic filter rule, or a combination of multiple atomic filter rules

| Field    | Type                                                                 | Required | Description                                                                                                          |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | yes      | node's field                                                                                                         |
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

**Note:**

- `enable_count`If this flag is true, this request is a get quantity. The remaining fields must be initialized, start is
  0, and limit is: 0, sort is "."
- `sort`If the caller does not specify it, the background specifies it as the container node ID by default.
- Paging parameters must be set, and the maximum query data at one time does not exceed 500.
- bk_cluster_id and cluster_uid cannot be empty or filled at the same time.

### Request Parameters Example

### Request Details Request Parameters

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "filter": {
    "condition": "OR",
    "rules": [
      {
        "field": "id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "bk_cluster_id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "hostname",
        "operator": "equal",
        "value": "name"
      }
    ]
  },
  "page": {
    "enable_count": false,
    "start": 0,
    "limit": 500
  }
}
```

### get quantity request parameters

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "filter": {
    "condition": "OR",
    "rules": [
      {
        "field": "id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "bk_cluster_id",
        "operator": "equal",
        "value": 10
      },
      {
        "field": "hostname",
        "operator": "equal",
        "value": "name"
      }
    ]
  },
  "page": {
    "enable_count": true,
    "start": 0,
    "limit": 0
  }
}
```

### Return Result Example

### Details interface response

```json
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "permission": null,
  "data": {
    "count": 0,
    "info": [
      {
        "name": "k8s",
        "roles": "master",
        "labels": {
          "env": "test"
        },
        "taints": {
          "type": "gpu"
        },
        "unschedulable": false,
        "internal_ip": [
          "127.0.0.1"
        ],
        "external_ip": [
          "127.0.0.1"
        ],
        "hostname": "name",
        "runtime_component": "runtime_component",
        "kube_proxy_mode": "ipvs",
        "pod_cidr": "127.0.0.128/26"
      }
    ]
  }
}
```

### kube node quantity interface response

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 1,
    "info": [
    ]
  },
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return Result Parameters Description

| Name       | Type   | Description                                                                        |
|------------|--------|------------------------------------------------------------------------------------|
| result     | bool   | Whether the request succeeded or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                         |
| message    | string | Error message returned by request failure                                          |
| permission | object | Permission information                                                             |
| data       | object | Data returned by request                                                           |
| request_id | string | Request chain id                                                                   |

#### data

| Field | Type  | Description       |
|-------|-------|-------------------|
| count | int   | Number of records |
| info  | array | Actual node data  |

#### info[x]

| Field             | Type   | Required | Description                                                                          |
|-------------------|--------|----------|--------------------------------------------------------------------------------------|
| name              | string | yes      | node name                                                                            |
| roles             | string | no       | node roles                                                                           |
| labels            | object | no       | label                                                                                |
| taints            | object | no       | taints                                                                               |
| unschedulable     | bool   | no       | Whether to turn off schedulable, true means not schedulable, false means schedulable |
| internal_ip       | array  | no       | internal ip                                                                          |
| external_ip       | array  | no       | external ip                                                                          |
| hostname          | string | no       | hostname                                                                             |
| runtime_component | string | no       | runtime components                                                                   |
| kube_proxy_mode   | string | no       | kube-proxy proxy mode                                                                |
| pod_cidr          | string | no       | The allocation range of the Pod address of this node                                 |

**Note:**

- If this request is to query details, count is 0. If the query is quantity, info is empty.

