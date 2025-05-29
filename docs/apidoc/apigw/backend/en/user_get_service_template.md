### Description

Get service template by service template ID.

### Parameters

| Name                | Type | Required | Description         |
|---------------------|------|----------|---------------------|
| service_template_id | int  | Yes      | Service template ID |

### Request Example

```json
{
  "service_template_id": 51
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
        "bk_biz_id": 3,
        "id": 51,
        "name": "mm2",
        "service_category_id": 12,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-05-26T09:46:15.259Z",
        "last_time": "2020-05-26T09:46:15.259Z",
        "bk_supplier_account": "0",
        "host_apply_enabled": false
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error    |
| message    | string | Error message returned in case of request failure                |
| permission | object | Permission information                                           |
| data       | object | Data returned by the request                                     |

#### data Field Description

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_biz_id           | int     | Business ID                                                |
| id                  | int     | Service template ID                                        |
| name                | array   | Service template name                                      |
| service_category_id | integer | Service category ID                                        |
| creator             | string  | Creator of the service template                            |
| modifier            | string  | Last modifier of the service template                      |
| create_time         | string  | Creation time                                              |
| last_time           | string  | Last update time                                           |
| bk_supplier_account | string  | Supplier account                                           |
| host_apply_enabled  | bool    | Whether to enable automatic application of host properties |
