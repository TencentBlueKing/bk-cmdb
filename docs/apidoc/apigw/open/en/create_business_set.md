### Description

Create a new business set (Version: v3.10.12+, Permission: Business Set Add Permission)

### Parameters

| Name            | Type   | Required | Description               |
|-----------------|--------|----------|---------------------------|
| bk_biz_set_attr | object | Yes      | Business set model fields |
| bk_scope        | object | Yes      | Selected business scope   |

#### bk_biz_set_attr

| Name              | Type   | Required | Description                          |
|-------------------|--------|----------|--------------------------------------|
| bk_biz_set_name   | string | Yes      | Business set name                    |
| bk_biz_maintainer | string | No       | Operations and Maintenance personnel |
| bk_biz_set_desc   | string | No       | Business set description             |

#### bk_scope

| Name      | Type   | Required | Description                        |
|-----------|--------|----------|------------------------------------|
| match_all | bool   | Yes      | Selected business scope flag       |
| filter    | object | No       | Selected business scope conditions |

#### filter

This parameter is a combination of filtering rules for business attribute fields, used to search for businesses based on
business attribute fields. The combination supports only AND operations, can be nested, and supports up to 2 levels.

| Name      | Type   | Required | Description                              |
|-----------|--------|----------|------------------------------------------|
| condition | string | Yes      | Rule operator                            |
| rules     | array  | Yes      | Selected business scope conditions rules |

#### rules

| Name     | Type   | Required | Description                                                        |
|----------|--------|----------|--------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                         |
| operator | string | Yes      | Operator. Optional values: equal, in                               |
| value    | -      | No       | Operand. Different operators correspond to different value formats |

**Note:**

- The input here is only for the required and system-built parameters for the `bk_biz_set_attr` parameter, and the rest
  of the parameters to be filled in depend on the user's own defined attribute fields.
- If the `match_all` field in `bk_scope` is true, it means that the selected business scope of the business set is all,
  and the parameter `filter` does not need to be filled in. If `match_all` field is false, `filter` needs to be
  non-empty, and users need to explicitly specify the selection range of the business.
- The business attribute enclosure type selected in the business set is organization and enumeration.

### Request Example

```python
{
    "bk_biz_set_attr":{
        "bk_biz_set_name":"biz_set",
        "bk_biz_set_desc":"xxx",
        "bk_biz_maintainer":"xxx"
    },
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
                    "field":"life_cycle",
                    "operator":"equal",
                    "value":1
                }
            ]
        }
    }
}
```

### Response Example

```python
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":5,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | int    | ID of the created business set                                              |
