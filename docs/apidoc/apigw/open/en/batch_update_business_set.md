### Description

Update Business Set Information (Version: v3.10.12+, Permission: Business Set Editing Permission)

### Parameters

| Name           | Type   | Required | Description              |
|----------------|--------|----------|--------------------------|
| bk_biz_set_ids | array  | Yes      | List of business set IDs |
| data           | object | Yes      | Business set data        |

#### data

| Name            | Type   | Required | Description               |
|-----------------|--------|----------|---------------------------|
| bk_biz_set_attr | object | No       | Business set model fields |
| bk_scope        | object | No       | Selected business scope   |

#### bk_biz_set_attr

| Name              | Type   | Required | Description              |
|-------------------|--------|----------|--------------------------|
| bk_biz_set_name   | string | Yes      | Business set name        |
| bk_biz_maintainer | string | No       | Operations personnel     |
| bk_biz_set_desc   | string | No       | Business set description |

#### bk_scope

| Name      | Type   | Required | Description                        |
|-----------|--------|----------|------------------------------------|
| match_all | bool   | Yes      | Selected business scope flag       |
| filter    | object | No       | Selected business scope conditions |

#### filter

This parameter is a combination of filtering rules for business property fields, used to search for hosts based on host
property fields. Combinations support only AND operations and can be nested, with a maximum of 2 layers.

| Name      | Type   | Required | Description                   |
|-----------|--------|----------|-------------------------------|
| condition | string | Yes      | Rule operator                 |
| rules     | array  | Yes      | Selected business scope rules |

#### rules

| Name     | Type   | Required | Description                                                        |
|----------|--------|----------|--------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                         |
| operator | string | Yes      | Operator. Optional values equal, in                                |
| value    | -      | No       | Operand. Different operators correspond to different value formats |

**Note:**

- The input parameters here only describe the required and system-built parameters. The rest of the parameters to be
  filled in depend on the attribute fields defined by the user.
- For batch scenarios (where the number of IDs in bk_biz_set_ids is greater than 1), changes to the `bk_biz_set_name`
  and `bk_scope` fields are not allowed.
- The maximum number of batch updates is 200.

### Request Example

```python
{
    "bk_biz_set_ids":[
        2
    ],
    "data":{
        "bk_biz_set_attr":{
            "bk_biz_set_name": "test",
            "bk_biz_set_desc":"xxx",
            "biz_set_maintainer":"xxx"
        },
        "bk_scope":{
            "match_all":false,
            "filter":{
                "condition":"AND",
                "rules":[
                    {
                        "field":"bk_sla",
                        "operator":"equal",
                        "value":"2"
                    }
                ]
            }
        }
    }
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {},
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
