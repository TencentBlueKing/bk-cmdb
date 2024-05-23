### Description

Query Module

### Parameters

| Name                | Type   | Required | Description                                                            |
|---------------------|--------|----------|------------------------------------------------------------------------|
| bk_supplier_account | string | No       | Developer account                                                      |
| bk_biz_id           | int    | Yes      | Business ID                                                            |
| bk_set_id           | int    | No       | Cluster ID                                                             |
| fields              | array  | Yes      | Query fields, fields come from the attributes defined in the module    |
| condition           | dict   | Yes      | Query condition, fields come from the attributes defined in the module |
| page                | dict   | Yes      | Paging condition                                                       |

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
        "bk_module_name",
        "bk_set_id"
    ],
    "condition": {
        "bk_module_name": "test"
    },
    "page": {
        "start": 0,
        "limit": 10
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
        "count": 2,
        "info": [
            {
                "bk_module_name": "test",
                "bk_set_id": 11,
                "default": 0
            },
            {
                "bk_module_name": "test",
                "bk_set_id": 12,
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

| Name  | Type  | Description                                                       |
|-------|-------|-------------------------------------------------------------------|
| count | int   | Number of data elements                                           |
| info  | array | Result set, where all fields are attributes defined in the module |

#### info

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_module_name      | string  | Module name                                                |
| bk_set_id           | int     | Cluster ID                                                 |
| default             | int     | Module type                                                |
| bk_bak_operator     | string  | Backup maintenance person                                  |
| bk_module_id        | int     | Model ID                                                   |
| bk_biz_id           | int     | Business ID                                                |
| bk_module_id        | int     | Module ID to which the host belongs                        |
| bk_module_type      | string  | Module type                                                |
| bk_parent_id        | int     | Parent node ID                                             |
| bk_supplier_account | string  | Developer account                                          |
| create_time         | string  | Creation time                                              |
| last_time           | string  | Update time                                                |
| host_apply_enabled  | bool    | Whether to enable automatic application of host properties |
| operator            | string  | Main maintainer                                            |
| service_category_id | integer | Service category ID                                        |
| service_template_id | int     | Service template ID                                        |
| set_template_id     | int     | Cluster template ID                                        |
| bk_created_at       | string  | Creation time                                              |
| bk_updated_at       | string  | Update time                                                |
| bk_created_by       | string  | Creator                                                    |

**Note: The return value here only describes the system's built-in property fields. The rest of the return value depends
on the user-defined property fields.**
