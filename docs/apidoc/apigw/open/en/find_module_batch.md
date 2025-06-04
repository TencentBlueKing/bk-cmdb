### Description

Batch obtain the attribute information of specified module instances under a specified business based on the business ID
and the list of module instance IDs, along with the desired module attribute list. (v3.8.6)

### Parameters

| Name      | Type  | Required | Description                                                                     |
|-----------|-------|----------|---------------------------------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                                                     |
| bk_ids    | array | Yes      | List of module instance IDs, i.e., bk_module_id list, up to 500                 |
| fields    | array | Yes      | Module attribute list, control which fields to return in the module information |

### Request Example

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        56,
        57,
        58,
        59,
        60
    ],
    "fields": [
        "bk_module_id",
        "bk_module_name",
        "create_time"
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [
        {
            "bk_module_id": 60,
            "bk_module_name": "sm1",
            "create_time": "2020-05-15T22:15:51.725+08:00",
            "default": 0
        },
        {
            "bk_module_id": 59,
            "bk_module_name": "m1",
            "create_time": "2020-05-12T21:04:47.286+08:00",
            "default": 0
        },
        {
            "bk_module_id": 58,
            "bk_module_name": "Pending recycle",
            "create_time": "2020-05-12T21:03:37.238+08:00",
            "default": 3
        },
        {
            "bk_module_id": 57,
            "bk_module_name": "Faulty machine",
            "create_time": "2020-05-12T21:03:37.183+08:00",
            "default": 2
        },
        {
            "bk_module_id": 56,
            "bk_module_name": "Idle machine",
            "create_time": "2020-05-12T21:03:37.122+08:00",
            "default": 1
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | array  | Data returned by the request                                       |

#### data Explanation

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_module_id        | int     | Module ID                                                  |
| bk_module_name      | string  | Module name                                                |
| default             | int     | Indicates the module type                                  |
| create_time         | string  | Creation time                                              |
| bk_set_id           | int     | Cluster ID                                                 |
| bk_bak_operator     | string  | Backup maintenance personnel                               |
| bk_biz_id           | int     | Business ID                                                |
| bk_module_type      | string  | Module type                                                |
| bk_parent_id        | int     | Parent node ID                                             |
| bk_supplier_account | string  | Developer account                                          |
| last_time           | string  | Update time                                                |
| host_apply_enabled  | bool    | Whether to enable automatic application of host properties |
| operator            | string  | Main maintainer                                            |
| service_category_id | integer | Service category ID                                        |
| service_template_id | int     | Service template ID                                        |
| set_template_id     | int     | Cluster template ID                                        |
| bk_created_at       | string  | Creation time                                              |
| bk_updated_at       | string  | Update time                                                |
| bk_created_by       | string  | Creator                                                    |

**Note: The returned value here only explains the built-in property fields. The rest of the returned values depend on
the user's own defined property fields.**
