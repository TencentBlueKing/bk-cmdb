### Functional description

search the business

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | No     | supplier account code |
| fields         |  array   | No     | need to show |
| condition      |  dict    | No     | search condition, legach field, please do not use this any more, use biz_property_filter instead |
| biz_property_filter    |  dict  | No     | business property filter |
| page           |  dict    | No     | page condition |

Note: a business has two status: normal or archived.
- search a archived business，please add rules `bk_data_status:disabled` to condition field.
- search a normal business，please do not add filed `bk_data_status` in condition , or add rule `bk_data_status: {"$ne":disabled"}` to condition.
- only one of `biz_property_filter` and `condition` parameters can take effect, and `condition` is not recommended to continue to use it.
- the number of array class elements involved in the parameter `biz_property_filter` shall not exceed 500.
  the number of `rules` involved in the parameter `biz_property_filter` does not exceed 20.
  the nesting level of parameter `biz_property_filter` shall not exceed 3 levels.

#### page

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| start    |  int    | Yes     | start record |
| limit    |  int    | Yes     | page limit, max is 200 |
| sort     |  string | No     | the field for sort |

### Request Parameters Example

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username": "xxx",
    "bk_token":"xxx",
    "bk_supplier_account":"123456789",
    "fields":[
        "bk_biz_id",
        "bk_biz_name"
    ],
    "biz_property_filter":{
        "condition":"AND",
        "rules":[
            {
                "field":"bk_biz_maintainer",
                "operator":"equal",
                "value":"admin"
            },
            {
                "condition":"OR",
                "rules":[
                    {
                        "field":"bk_biz_name",
                        "operator":"in",
                        "value":[
                            "test"
                        ]
                    },
                    {
                        "field":"bk_biz_id",
                        "operator":"equal",
                        "value":0
                    }
                ]
            }
        ]
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":""
    }
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |
| permission    | object | permission Information    |
| request_id    | string | request chain id    |

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     | int       | the num of record |
| info      | array     | business info |

#### info

| Field | Type | Description |
|-----------|-----------|-----------|
| bk_biz_id | int | business id |
| bk_biz_name | string | business name |
|default | int | indicates the type of business |
