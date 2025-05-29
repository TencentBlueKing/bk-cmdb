### Description

Query Cluster

### Parameters

| Name                | Type   | Required | Description                                                                                                    |
|---------------------|--------|----------|----------------------------------------------------------------------------------------------------------------|
| bk_supplier_account | string | No       | Developer account                                                                                              |
| bk_biz_id           | int    | Yes      | Business ID                                                                                                    |
| fields              | array  | Yes      | Query fields, all fields are attributes defined in the set, including preset fields and user-defined fields    |
| condition           | dict   | Yes      | Query condition, all fields are attributes defined in the set, including preset fields and user-defined fields |
| filter | object | No       |  set condition range                                 |
| time_condition | object | No       |  set time range                                      |
| page                | dict   | Yes      | Paging condition                                                                                               |

- Only one of the parameters `filter` and `condition` can be effective, and it is not recommended to continue using the parameter `condition`.
- The number of array elements involved in the parameter `filter` does not exceed 500. The number of `rules` involved in the parameter `filter` does not exceed 20. The nesting level of the parameter `filter` does not exceed 3.

#### filter

| Field     | Type   | Required | Description                                    |
| --------- | ------ | -------- | ---------------------------------------------- |
| condition | string | Yes      | Rule operator                                  |
| rules     | array  | Yes      | Filtering rules for the scope of business sets |

#### rules

Filter rules are triplets `field`, `operator`, `value`

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | Yes      | Field name                                                   |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats |

Assembly rules can refer to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### time_condition

| Field | Type   | Required | Description                           |
| ----- | ------ | -------- | ------------------------------------- |
| oper  | string | Yes      | Operator, currently only supports and |
| rules | array  | Yes      | Time query conditions                 |

#### time_condition.rules

| Field | Type   | Required | Description                               |
| ----- | ------ | -------- | ----------------------------------------- |
| field | string | Yes      | Takes the value of the model's field name |
| start | string | Yes      | Start time, format: yyyy-MM-dd hh:mm:ss   |
| end   | string | Yes      | End time, format: yyyy-MM-dd hh:mm:ss     |

#### page

| Name  | Type   | Required | Description           |
|-------|--------|----------|-----------------------|
| start | int    | Yes      | Record start position |
| limit | int    | Yes      | Limit per page        |
| sort  | string | No       | Sorting field         |

### Request Example

```json
{
    "bk_biz_id": 2,
    "fields": [
        "bk_set_name"
    ],
    "filter": {
      "condition": "AND",
      "rules": [
        {
          "field": "bk_set_name",
          "operator": "equal",
          "value": "test"
        }
      ]
    },
    "time_condition": {
      "oper": "and",
      "rules": [
        {
          "field": "create_time",
          "start": "2021-05-13 01:00:00",
          "end": "2021-05-14 01:00:00"
        }
      ]
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_set_name"
    }
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "count": 1,
        "info": [
            {
                "bk_set_name": "test",
                "default": 0
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |

#### data

| Name  | Type  | Description                                                    |
|-------|-------|----------------------------------------------------------------|
| count | int   | Number of data elements                                        |
| info  | array | Result set, where all fields are attributes defined in the set |

#### info

| Name                 | Type   | Description                                                |
|----------------------|--------|------------------------------------------------------------|
| bk_set_name          | string | Cluster name                                               |
| default              | int    | 0-normal cluster, 1-built-in module set, default is 0      |
| bk_biz_id            | int    | Business ID                                                |
| bk_capacity          | int    | Design capacity                                            |
| bk_parent_id         | int    | Parent node ID                                             |
| bk_set_id            | int    | Cluster ID                                                 |
| bk_service_status    | string | Service status: 1/2(1: open, 2: closed)                    |
| bk_set_desc          | string | Cluster description                                        |
| bk_set_env           | string | Environment type: 1/2/3(1: test, 2: experience, 3: formal) |
| create_time          | string | Creation time                                              |
| last_time            | string | Update time                                                |
| bk_supplier_account  | string | Developer account                                          |
| description          | string | Description of the data                                    |
| set_template_version | array  | Current version of cluster template                        |
| set_template_id      | int    | Cluster template ID                                        |
| bk_created_at        | string | Creation time                                              |
| bk_updated_at        | string | Update time                                                |
| bk_created_by        | string | Creator                                                    |

**Note: The return value here only describes the system's built-in property fields. The rest of the return value depends
on the user-defined property fields.**
