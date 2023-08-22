### Functional description

list container clusters (v3.12.1+, permissions: biz access)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                                                |
|-----------|--------|----------|--------------------------------------------------------------------------------------------|
| bk_biz_id | int    | yes      | business ID                                                                                |
| filter    | object | no       | container cluster query scope                                                              |
| fields    | array  | no       | the container cluster attribute to be queried, if not written, it means to search all data |
| page      | object | yes      | paging condition                                                                           |

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
| field    | string                                                               | yes      | cluster's field                                                                                                      |
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
- `sort`If the caller does not specify it, the background specifies it as the container cluster ID by default.
- Paging parameters must be set, and the maximum query data at one time does not exceed 500.

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
    "condition": "AND",
    "rules": [
      {
        "field": "scheduling_engine",
        "operator": "equal",
        "value": "k8s"
      },
      {
        "field": "version",
        "operator": "equal",
        "value": "1.1.0"
      }
    ]
  },
  "page": {
    "start": 0,
    "limit": 500,
    "enable_count": false
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
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "scheduling_engine",
        "operator": "equal",
        "value": "k8s"
      },
      {
        "field": "version",
        "operator": "equal",
        "value": "1.1.0"
      }
    ]
  },
  "page": {
    "start": 0,
    "limit": 0,
    "enable_count": true
  }
}
```

### Return Result Example

### Details interface response

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
        "name": "cluster",
        "scheduling_engine": "k8s",
        "uid": "xxx",
        "xid": "xxx",
        "version": "1.1.0",
        "network_type": "underlay",
        "region": "xxx",
        "vpc": "xxx",
        "network": "127.0.0.0/21",
        "type": "INDEPENDENT_CLUSTER"
      }
    ]
  },
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### kube cluster quantity interface response

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

| Field | Type  | Description         |
|-------|-------|---------------------|
| count | int   | Number of records   |
| info  | array | Actual cluster data |

#### info[x]

| Field             | Type   | Required | Description                             |
|-------------------|--------|----------|-----------------------------------------|
| name              | string | no       | cluster                                 |
| scheduling_engine | string | no       | scheduling engine                       |
| uid               | string | no       | cluster own ID                          |
| xid               | string | no       | associated cluster ID                   |
| version           | string | no       | cluster version                         |
| network_type      | string | no       | network type                            |
| region            | string | no       | the region where the cluster is located |
| vpc               | string | no       | vpc network                             |
| network           | array  | no       | cluster network                         |
| type              | string | no       | cluster type                            |

**Note:**

- If this request is to query details, count is 0. If the query is quantity, info is empty.
