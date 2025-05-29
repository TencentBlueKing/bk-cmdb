### Description

Query Business (Permission: Business Query Permission)

### Parameters

| Name                | Type   | Required | Description                                                                                                                                                                               |
|---------------------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_supplier_account | string | No       | Developer account                                                                                                                                                                         |
| fields              | array  | No       | Specify the fields to query. If not filled in, the system will return all fields of the business                                                                                          |
| condition           | dict   | No       | Query conditions, parameters for any property of the business. If not written, it means to search all data. (Legacy field, please do not continue to use, please use biz_property_filter) |
| biz_property_filter | object | No       | Business attribute combination query conditions                                                                                                                                           |
| time_condition | object | No | Query criteria for querying business by time |
| page                | dict   | No       | Paging conditions                                                                                                                                                                         |

Note: Businesses are divided into two types, non-archived businesses and archived businesses.

- To query archived businesses, add the condition `bk_data_status:disabled` in the condition.
- To query non-archived businesses, do not include the field "bk_data_status" or add the
  condition `bk_data_status: {"$ne": "disabled"}` in the condition.
- Only one of the parameters `biz_property_filter` and `condition` can be effective, and it is not recommended to
  continue using the parameter `condition`.
- The number of array elements involved in the parameter `biz_property_filter` does not exceed 500. The number
  of `rules` involved in the parameter `biz_property_filter` does not exceed 20. The nesting level of the
  parameter `biz_property_filter` does not exceed 3.

#### biz_property_filter

| Name      | Type   | Required | Description                         |
|-----------|--------|----------|-------------------------------------|
| condition | string | Yes      | Aggregation condition               |
| rules     | array  | Yes      | Rules for the aggregation condition |

#### rules

| Name     | Type   | Required | Description |
|----------|--------|----------|-------------|
| field    | string | Yes      | Field       |
| operator | string | Yes      | Operator    |
| value    | object | Yes      | Value       |

#### time_condition

| Field   | Type   | Required| Description              |
|-------|--------|-----|--------------------|
| oper  | string |Yes| Operator, currently only and is supported|
| rules | array  |Yes| Time query criteria         |

#### page

| Name  | Type   | Required | Description                                                                                                              |
|-------|--------|----------|--------------------------------------------------------------------------------------------------------------------------|
| start | int    | Yes      | Record start position                                                                                                    |
| limit | int    | Yes      | Limit per page, maximum 200                                                                                              |
| sort  | string | No       | Sorting field. Adding "-" in front of the field, such as sort: "-field", represents sorting by field in descending order |

### Request Example

```json
{
    "fields": ["bk_biz_id", "bk_biz_name"],
    "biz_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_biz_maintainer",
                "operator": "equal",
                "value": "admin"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_biz_name",
                        "operator": "in",
                        "value": ["test"]
                    },
                    {
                        "field": "bk_biz_id",
                        "operator": "equal",
                        "value": 1
                    }
                ]
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
        "sort": ""
    }
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {
        "count": 1,
        "info": [
            {
                "bk_biz_id": 1,
                "bk_biz_name": "esb-test",
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

| Name  | Type  | Description                 |
|-------|-------|-----------------------------|
| count | int   | Number of records           |
| info  | array | Actual data of the business |

#### info

| Name                | Type   | Description                                |
|---------------------|--------|--------------------------------------------|
| bk_biz_id           | int    | Business ID                                |
| bk_biz_name         | string | Business name                              |
| bk_biz_maintainer   | string | Operation and maintenance personnel        |
| bk_biz_productor    | string | Product personnel                          |
| bk_biz_developer    | string | Developer                                  |
| bk_biz_tester       | string | Tester                                     |
| time_zone           | string | Time zone                                  |
| language            | string | Language, "1" for Chinese, "2" for English |
| bk_supplier_account | string | Developer account                          |
| create_time         | string | Creation time                              |
| last_time           | string | Update time                                |
| default             | int    | Business type                              |
| operator            | string | Main maintainer                            |
| life_cycle          | string | Business life cycle                        |
| bk_created_at       | string | Creation time                              |
| bk_updated_at       | string | Update time                                |
| bk_created_by       | string | Creator                                    |

**Note: The return value here only describes the system's built-in property fields. The rest of the return value depends
on the user-defined property fields.**
