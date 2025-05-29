### Description

Update service template information (Permission: Service Template Editing Permission)

### Parameters

| Name                | Type   | Required                                                               | Description           |
|---------------------|--------|------------------------------------------------------------------------|-----------------------|
| name                | string | Either `service_category_id` or `name` is required, both can be filled | Service template name |
| service_category_id | int    | Either `service_category_id` or `name` is required, both can be filled | Service category ID   |
| id                  | int    | Yes                                                                    | Service template ID   |
| bk_biz_id           | int    | Yes                                                                    | Business ID           |

### Request Example

```python
{
  "bk_biz_id": 1,
  "name": "test1",
  "id": 50,
  "service_category_id": 3
}
```

### Response Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "bk_biz_id": 1,
    "id": 50,
    "name": "test1",
    "service_category_id": 3,
    "creator": "admin",
    "modifier": "admin",
    "create_time": "2019-06-05T11:22:22.951+08:00",
    "last_time": "2019-06-05T11:22:22.951+08:00",
    "bk_supplier_account": "0",
    "host_apply_enabled": false
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Updated service template information                                |

#### data

| Name                | Type   | Description                                           |
|---------------------|--------|-------------------------------------------------------|
| id                  | int    | Service template ID                                   |
| name                | string | Service template name                                 |
| bk_biz_id           | int    | Business ID                                           |
| service_category_id | int    | Service category ID                                   |
| creator             | string | Creator                                               |
| modifier            | string | Last modifier                                         |
| create_time         | string | Creation time                                         |
| last_time           | string | Last update time                                      |
| bk_supplier_account | string | Supplier account                                      |
| host_apply_enabled  | bool   | Whether to enable host property automatic application |
