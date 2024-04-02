### Function Description

Query business in the centralized business (Version: v3.10.12+, Permission: Business set access permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field         | Type   | Required | Description                                                  |
| ------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_set_id | int    | Yes      | Business set ID                                              |
| filter        | object | No       | Business attribute combination query conditions              |
| fields        | array  | No       | Specify the fields to query, parameters can be any business attribute. If no field information is provided, the system will return all fields of the business |
| page          | object | Yes      | Pagination conditions                                        |

#### filter

Query conditions. Supports combination of AND and OR. Can be nested, with a maximum nesting level of 2.

| Field     | Type   | Required | Description                               |
| --------- | ------ | -------- | ----------------------------------------- |
| condition | string | Yes      | Rule operator                             |
| rules     | array  | Yes      | Filtering rules for the scope of business |

#### rules

Filter rules are triplets `field`, `operator`, `value`

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | Yes      | Field name                                                   |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats |

Assembly rules can refer to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### page

| Field        | Type   | Required | Description                                                  |
| ------------ | ------ | -------- | ------------------------------------------------------------ |
| start        | int    | Yes      | Record start position                                        |
| limit        | int    | Yes      | Number of records per page, maximum 500                      |
| enable_count | bool   | Yes      | Flag to get the count of query objects                       |
| sort         | string | No       | Sorting field, by adding a "-", e.g., sort: ""-field"", it represents sorting in descending order by the field |

**Note:**

- If `enable_count` is set to true, it indicates that this request is to obtain the count. In this case, other fields must have initial values, start is 0, limit is 0, sort is "".

### Request Parameters Example

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
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

### Detailed Information Response Example

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
    "request_id": "dsda1122adasadadada2222"
}
```

### Query Business Count Response Example

```
pythonCopy code{
    "result":true,
    "code":0,
    "message":"",
    "permission":null,
    "data":{
        "count":10,
        "info":[

        ]
    },
    "request_id": "dsda1122adasadadada2222"
}
```

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Success or failure of the request. true: success; false: failure |
| code       | int    | Error code. 0 represents success, >0 represents failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| data       | object | Data returned by the request                                 |
| request_id | string | Request chain ID                                             |

#### data

| Field | Type  | Description          |
| ----- | ----- | -------------------- |
| count | int   | Number of records    |
| info  | array | Actual business data |

#### info

| Field               | Type   | Description                                |
| ------------------- | ------ | ------------------------------------------ |
| bk_biz_id           | int    | Business ID                                |
| bk_biz_name         | string | Business name                              |
| bk_biz_maintainer   | string | Operations personnel                       |
| bk_biz_productor    | string | Product personnel                          |
| bk_biz_developer    | string | Development personnel                      |
| bk_biz_tester       | string | Testing personnel                          |
| time_zone           | string | Time zone                                  |
| language            | string | Language, "1" for Chinese, "2" for English |
| bk_supplier_account | string | Supplier account                           |
| create_time         | string | Creation time                              |
| last_time           | string | Update time                                |
| default             | int    | Indicates business type                    |
| operator            | string | Main maintainer                            |
| life_cycle          | string | Business lifecycle                         |
| bk_created_at       | string | Creation time                              |
| bk_updated_at       | string | Update time                                |
| bk_created_by       | string | Creator                                    |

**Note:**

- The returned values here only explain the system's built-in attribute fields, other returned values depend on user-defined attribute fields.