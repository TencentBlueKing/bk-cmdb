### Description

Query container clusters (v3.12.1+, Permission: Business access)

### Parameters

| Name      | Type   | Required | Description                                                                             |
|-----------|--------|----------|-----------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                             |
| filter    | object | No       | Container cluster query scope                                                           |
| fields    | array  | No       | Container cluster properties to be queried. If not specified, all data will be searched |
| page      | object | Yes      | Pagination conditions                                                                   |

#### filter

This parameter is a combination of container cluster property field filtering rules, used to search for container
clusters based on container cluster property fields. The combination supports AND and OR two ways, and allows nesting,
with a maximum nesting of 2 levels.

| Name      | Type   | Required | Description                               |
|-----------|--------|----------|-------------------------------------------|
| condition | string | Yes      | Rule operator                             |
| rules     | array  | Yes      | Filtering rules for the range of clusters |

#### rules

The filtering rule is a triple `field`, `operator`, `value`.

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can be referred
to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### page

| Name         | Type   | Required | Description                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------|
| start        | int    | Yes      | Record start position                                                                                              |
| limit        | int    | Yes      | Number of records per page, maximum 500                                                                            |
| enable_count | bool   | Yes      | Flag for whether this request is for obtaining the quantity or details                                             |
| sort         | string | No       | Sorting field, by adding - in front of the field, such as sort:"-field", it can indicate descending order by field |

**Note:**

- `enable_count` If this flag is true, it means this request is to obtain the quantity. At this time, other fields must
  be initialized, start is 0, limit is 0, sort is "".
- If `sort` is not specified by the caller, the backend defaults to the container cluster ID.
- Pagination parameters must be set, and the maximum number of queried data at a time should not exceed 500.

### Request Example

#### Detailed Information Request Parameters

```json
{
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

#### Quantity Request Example

```json
{
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

### Response Example

#### Detailed Information Interface Response

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
        "type": "INDEPENDENT_CLUSTER",
        "environment": "xxx",
        "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
        "bk_project_name": "test",
        "bk_project_code": "test"
      }
    ]
  },
}
```

#### Quantity Interface Response

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

| Name  | Type  | Description         |
|-------|-------|---------------------|
| count | int   | Number of records   |
| info  | array | Actual cluster data |

#### info[x]

| Name              | Type   | Description           |
|-------------------|--------|-----------------------|
| name              | string | Cluster name          |
| scheduling_engine | string | Scheduling engine     |
| uid               | string | Cluster's own ID      |
| xid               | string | Associated cluster ID |
| version           | string | Cluster version       |
| network_type      | string | Network type          |
| region            | string | Region                |
| vpc               | string | VPC network           |
| network           | array  | Cluster network       |
| type              | string | Cluster type          |
| environment       | string | Environment           |
| bk_project_id     | string | Project ID            |
| bk_project_name   | string | Project name          |
| bk_project_code   | string | Project English name  |

**Note:**

- If this request is to query detailed information, count is 0. If it is to query the quantity, info is empty.
