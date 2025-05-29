### Description

Query business sets (Version: v3.10.12+, Permission: Business set view permission)

### Parameters

| Name              | Type   | Required | Description                                                                                  |
|-------------------|--------|----------|----------------------------------------------------------------------------------------------|
| bk_biz_set_filter | object | No       | Business set condition range                                                                 |
| time_condition    | object | No       | Business set time range                                                                      |
| fields            | array  | No       | Query conditions, parameters can be any business attribute. If not provided, search all data |
| page              | object | Yes      | Pagination conditions                                                                        |

#### bk_biz_set_filter

This parameter is a combination of filtering rules for business set attribute fields, used to search for business sets
based on business set attribute fields. The combination supports AND and OR, allows nesting, with a maximum nesting
level of 2.

| Name      | Type   | Required | Description                                    |
|-----------|--------|----------|------------------------------------------------|
| condition | string | Yes      | Rule operator                                  |
| rules     | array  | Yes      | Filtering rules for the scope of business sets |

#### rules

Filter rules are triplets `field`, `operator`, `value`

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can refer
to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### time_condition

| Name  | Type   | Required | Description                           |
|-------|--------|----------|---------------------------------------|
| oper  | string | Yes      | Operator, currently only supports and |
| rules | array  | Yes      | Time query conditions                 |

#### rules

| Name  | Type   | Required | Description                               |
|-------|--------|----------|-------------------------------------------|
| field | string | Yes      | Takes the value of the model's field name |
| start | string | Yes      | Start time, format: yyyy-MM-dd hh:mm:ss   |
| end   | string | Yes      | End time, format: yyyy-MM-dd hh:mm:ss     |

#### page

| Name         | Type   | Required | Description                                                                                                    |
|--------------|--------|----------|----------------------------------------------------------------------------------------------------------------|
| start        | int    | Yes      | Record start position                                                                                          |
| limit        | int    | Yes      | Number of records per page, maximum 500                                                                        |
| enable_count | bool   | Yes      | Flag to get the count of query objects                                                                         |
| sort         | string | No       | Sorting field, by adding a "-", e.g., sort: ""-field"", it represents sorting in descending order by the field |

**Note:**

- If `enable_count` is set to true, it indicates that this request is to obtain the count. In this case, other fields
  must have initial values, start is 0, limit is 0, sort is "".
- If `sort` is not specified by the caller, the backend defaults it to the business set ID.

### Request Example

```python
{
    "bk_biz_set_filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"bk_biz_set_id",
                "operator":"equal",
                "value":10
            },
            {
                "field":"bk_biz_maintainer",
                "operator":"equal",
                "value":"admin"
            }
        ]
    },
    "time_condition":{
        "oper":"and",
        "rules":[
            {
                "field":"create_time",
                "start":"2021-05-13 01:00:00",
                "end":"2021-05-14 01:00:00"
            }
        ]
    },
    "fields": [
        "bk_biz_id"
    ],
    "page":{
        "start":0,
        "limit":500,
        "enable_count":false,
        "sort":"bk_biz_set_id"
    }
}
```

### Response Example

#### Detailed Information Interface Response

```python
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":0,
        "info":[
            {
                "bk_biz_set_id":10,
                "bk_biz_set_name":"biz_set",
                "bk_biz_set_desc":"dba",
                "bk_biz_maintainer":"tom",
                "create_time":"2021-09-06T08:10:50.168Z",
                "last_time":"2021-10-15T02:30:01.867Z",
                "bk_scope":{
                    "match_all":true
                }
            },
            {
                "bk_biz_set_id":11,
                "bk_biz_set_name":"biz_set1",
                "bk_biz_set_desc":"dba",
                "bk_biz_maintainer":"tom",
                "create_time":"2021-09-06T08:10:50.168Z",
                "last_time":"2021-10-15T02:30:01.867Z",
                "bk_scope":{
                    "match_all":false,
                    "filter":{
                        "condition":"AND",
                        "rules":[
                            {
                                "field":"bk_sla",
                                "operator":"equal",
                                "value":"3"
                            },
                            {
                                "field":"bk_biz_maintainer",
                                "operator":"equal",
                                "value":"admin"
                            }
                        ]
                    }
                }
            }
        ]
    },
}
```

#### Business Set Count Interface Response

```
pythonCopy code{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":2,
        "info":[
        ]
    },
}
```

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Success or failure of the request. true: success; false: failure |
| code       | int    | Error code. 0 represents success, >0 represents failure error    |
| message    | string | Error message returned in case of failure                        |
| permission | object | Permission information                                           |
| data       | object | Data returned by the request                                     |

#### data

| Name  | Type  | Description          |
|-------|-------|----------------------|
| count | int   | Number of records    |
| info  | array | Actual business data |

#### info

| Name              | Type   | Description                    |
|-------------------|--------|--------------------------------|
| bk_biz_set_id     | int    | Business set ID                |
| create_time       | string | Business set creation time     |
| last_time         | string | Business set modification time |
| bk_biz_set_name   | string | Business set name              |
| bk_biz_maintainer | string | Operations personnel           |
| bk_biz_set_desc   | string | Business set description       |
| bk_scope          | object | Selected business scope        |
| bk_created_at     | string | Creation time                  |
| bk_created_by     | string | Creator                        |
| bk_updated_at     | string | Update time                    |

#### bk_scope

| Name      | Type   | Description                        |
|-----------|--------|------------------------------------|
| match_all | bool   | Selected business scope flag       |
| filter    | object | Selected business range conditions |

#### filter

This parameter is a combination of filtering rules for business attributes, used to search for businesses based on
business attributes. The combination only supports AND operations, can be nested, with a maximum nesting level of 2.

| Name      | Type   | Description                               |
|-----------|--------|-------------------------------------------|
| condition | string | Rule operator                             |
| rules     | array  | Filtering rules for the scope of business |

#### rules

| Name     | Type   | Description                                                        |
|----------|--------|--------------------------------------------------------------------|
| field    | string | Field name                                                         |
| operator | string | Operator, optional values are equal, in                            |
| value    | -      | Operand, different operators correspond to different value formats |

**Note:**

- If this request is to query detailed information, count is 0. If querying for the count, info is empty.
- The input here for the `info` parameter only explains the required and system built-in parameters, other parameters to
  be filled depend on user-defined attribute fields.
