### Description

Query container nodes (v3.12.1+, Permission: Business access)

### Parameters

| Name      | Type   | Required | Description                                                                           |
|-----------|--------|----------|---------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                           |
| filter    | object | No       | Container node query scope                                                            |
| fields    | array  | No       | Container node properties to be queried. If not specified, all data will be searched. |
| page      | object | Yes      | Pagination information                                                                |

#### filter

This parameter is a combination of filter rules for container node properties, used to search for container clusters
based on container node properties. Combinations support both AND and OR, allowing nesting, with a maximum nesting of 2
levels.

| Name      | Type   | Required | Description                            |
|-----------|--------|----------|----------------------------------------|
| condition | string | Yes      | Rule operator                          |
| rules     | array  | Yes      | Filtering rules for the range of nodes |

#### rules

The filtering rule is a triple `field`, `operator`, `value`

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can be referred
to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### page

| Name         | Type   | Required | Description                                                                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | Yes      | Record start position                                                                                                                                              |
| limit        | int    | Yes      | Number of records per page, maximum 500                                                                                                                            |
| enable_count | bool   | Yes      | Flag for whether this request is for obtaining the quantity or details. If this flag is true, other fields must be initialized. start is 0, limit is 0, sort is "" |
| sort         | string | No       | Sorting field, adding a `-` before the field, such as sort: ""-field"", represents sorting the field in descending order                                           |

**Note:**

- `enable_count` If this flag is true, it means this request is to obtain the quantity. At this time, other fields must
  be initialized, start is 0, limit is 0, sort is "".
- If `sort` is not specified by the caller, the backend defaults to the node ID.
- Pagination parameters must be set, and the maximum number of queried data at a time should not exceed 500.

### Request Example

#### Detailed Information Request Parameter

```json
{
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

#### Quantity Request Parameter

```json
{
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

### Response Example

#### Detailed Information Interface Response

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

#### Quantity Response Example

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

| Name  | Type  | Description                       |
|-------|-------|-----------------------------------|
| count | int   | Number of records                 |
| info  | array | Actual data, details of the nodes |

#### info[x]

| Name              | Type   | Description                                                                       |
|-------------------|--------|-----------------------------------------------------------------------------------|
| name              | string | Node name                                                                         |
| roles             | string | Node type                                                                         |
| labels            | object | Labels                                                                            |
| taints            | object | Taints                                                                            |
| unschedulable     | bool   | Whether scheduling is closed, true means not schedulable, false means schedulable |
| internal_ip       | array  | Internal IP                                                                       |
| external_ip       | array  | External IP                                                                       |
| hostname          | string | Hostname                                                                          |
| runtime_component | string | Runtime component                                                                 |
| kube_proxy_mode   | string | kube-proxy proxy mode                                                             |
| pod_cidr          | string | The allocation range of Pod addresses on this node                                |

**Note:**

- If this request is to query detailed information, count is 0. If the query is for quantity, info is empty.
