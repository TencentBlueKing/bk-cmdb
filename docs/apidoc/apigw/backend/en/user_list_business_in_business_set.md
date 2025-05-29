### Description

Query business in the centralized business (Version: v3.10.12+, Permission: Business set access permission)

### Parameters

| Name          | Type   | Required | Description                                                                                                                                                   |
|---------------|--------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_set_id | int    | Yes      | Business set ID                                                                                                                                               |
| filter        | object | No       | Business attribute combination query conditions                                                                                                               |
| fields        | array  | No       | Specify the fields to query, parameters can be any business attribute. If no field information is provided, the system will return all fields of the business |
| page          | object | Yes      | Pagination conditions                                                                                                                                         |

#### filter

Query conditions. Supports combination of AND and OR. Can be nested, with a maximum nesting level of 2.

| Name      | Type   | Required | Description                               |
|-----------|--------|----------|-------------------------------------------|
| condition | string | Yes      | Rule operator                             |
| rules     | array  | Yes      | Filtering rules for the scope of business |

#### rules

Filter rules are triplets `field`, `operator`, `value`

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can refer
to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

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

### Request Example

```python
{
    "bk_biz_set_id":2,
    "filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"xxx",
                "operator":"equal",
                "value":"xxx"
            },
            {
                "field":"xxx",
                "operator":"in",
                "value":[
                    "xxx"
                ]
            }
        ]
    },
    "fields":[
        "bk_biz_id",
        "bk_biz_name"
    ],
    "page":{
        "start":0,
        "limit":10,
        "enable_count":false,
        "sort":"bk_biz_id"
    }
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission":null,
    "data": {
        "count": 0,
        "info": [
            {
                "bk_biz_id": 1,
                "bk_biz_name": "xxx"
            }
        ]
    },
}
```

### Response Parameters
