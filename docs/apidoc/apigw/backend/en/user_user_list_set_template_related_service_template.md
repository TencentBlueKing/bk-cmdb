### Description

Retrieve a list of service templates for a specified business and cluster template ID.

### Parameters

| Name            | Type | Required | Description         |
|-----------------|------|----------|---------------------|
| set_template_id | int  | Yes      | Cluster template ID |
| bk_biz_id       | int  | Yes      | Business ID         |

### Request Example

```json
{
  "set_template_id": 1,
  "bk_biz_id": 3
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
            "bk_biz_id": 3,
            "id": 48,
            "name": "sm1",
            "service_category_id": 2,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:14:57.691Z",
            "last_time": "2020-05-15T14:14:57.691Z",
            "bk_supplier_account": "0",
            "host_apply_enabled": false
        },
        {
            "bk_biz_id": 3,
            "id": 49,
            "name": "sm2",
            "": 16,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:19:09.813Z",
            "last_time": "2020-05-15T14:19:09.813Z",
            "bk_supplier_account": "0",
            "host_apply_enabled": false
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | array  | Data returned by the request                                       |

#### data

| Name                | Type   | Description                                                |
|---------------------|--------|------------------------------------------------------------|
| bk_biz_id           | int    | Business ID                                                |
| id                  | int    | Service template ID                                        |
| name                | string | Service template name                                      |
| service_category_id | int    | Service category ID                                        |
| creator             | string | Creator                                                    |
| modifier            | string | Last modifier                                              |
| create_time         | string | Creation time                                              |
| last_time           | string | Update time                                                |
| bk_supplier_account | string | Supplier account                                           |
| host_apply_enabled  | bool   | Whether to enable automatic application of host properties |
