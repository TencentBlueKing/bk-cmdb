### Description

Query modules under a business based on conditions (v3.9.7)

### Parameters

| Name                    | Type   | Required | Description                                                                     |
|-------------------------|--------|----------|---------------------------------------------------------------------------------|
| bk_biz_id               | int    | Yes      | Business ID                                                                     |
| bk_set_ids              | array  | No       | Cluster ID list, up to 200                                                      |
| bk_service_template_ids | array  | No       | Service template ID list                                                        |
| fields                  | array  | Yes      | Module attribute list, control which fields to return in the module information |
| page                    | object | Yes      | Pagination information                                                          |

#### page Field Explanation

| Name  | Type | Required | Description                                |
|-------|------|----------|--------------------------------------------|
| start | int  | Yes      | Record start position                      |
| limit | int  | Yes      | Number of records per page, maximum is 500 |

### Request Example

```json
{
    "bk_biz_id": 2,
    "bk_set_ids":[1,2],
    "bk_service_template_ids": [3,4],
    "fields":["bk_module_id", "bk_module_name"],
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
                "bk_module_id": 8,
                "bk_module_name": "license"
            },
            {
                "bk_module_id": 12,
                "bk_module_name": "gse_proc"
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### Explanation of data field:

| Name  | Type         | Description            |
|-------|--------------|------------------------|
| count | int          | Number of records      |
| info  | object array | Actual data of modules |

#### Explanation of data.info field:

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_module_id        | int     | Module ID                                                  |
| bk_module_name      | string  | Module name                                                |
| default             | int     | Indicates the module type                                  |
| create_time         | string  | Creation time                                              |
| bk_set_id           | int     | Cluster ID                                                 |
| bk_bak_operator     | string  | Backup maintainer                                          |
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
