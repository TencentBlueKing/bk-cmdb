### Description

Retrieve a list of service templates based on the business ID, with an option to further filter by service category ID.

### Parameters

| Name                 | Type      | Required | Description                                                                                                                                  |
|----------------------|-----------|----------|----------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id            | int       | Yes      | Business ID                                                                                                                                  |
| service_category_id  | int       | No       | Service category ID                                                                                                                          |
| search               | string    | No       | Search by service template name; default is empty                                                                                            |
| is_exact             | bool      | No       | Whether to match the service template name exactly; default is false. Effective when used in conjunction with the search parameter (v3.9.19) |
| service_template_ids | int array | No       | Service template IDs                                                                                                                         |
| page                 | object    | Yes      | Pagination parameters                                                                                                                        |

#### page

| Name  | Type   | Required | Description                             |
|-------|--------|----------|-----------------------------------------|
| start | int    | Yes      | Record start position                   |
| limit | int    | Yes      | Number of records per page, maximum 500 |
| sort  | string | No       | Sorting field                           |

### Request Example

```json
{
    "bk_biz_id": 1,
    "service_category_id": 1,
    "service_template_ids":[5,6],
    "search": "test2",
    "is_exact": true,
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-name"
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
        "count": 1,
        "info": [
            {
                "bk_biz_id": 1,
                "id": 50,
                "name": "test2",
                "service_category_id": 1,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2019-09-18T20:31:29.607+08:00",
                "last_time": "2019-09-18T20:31:29.607+08:00",
                "bk_supplier_account": "0",
                "host_apply_enabled": false
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data Field Explanation

| Name  | Type  | Description              |
|-------|-------|--------------------------|
| count | int   | Total number of records  |
| info  | array | List of returned results |

#### info Field Explanation

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_biz_id           | int     | Business ID                                                |
| id                  | int     | Service template ID                                        |
| name                | array   | Service template name                                      |
| service_category_id | integer | Service category ID                                        |
| creator             | string  | Creator of this data                                       |
| modifier            | string  | Last modifier of this data                                 |
| create_time         | string  | Creation time                                              |
| last_time           | string  | Update time                                                |
| bk_supplier_account | string  | Supplier account                                           |
| host_apply_enabled  | bool    | Whether to enable automatic application of host properties |
