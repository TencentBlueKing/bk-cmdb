### Description

Create a cluster template with the specified name under the specified business ID. The created cluster template includes
service templates based on the specified service template IDs. (Permission: Cluster template creation permission)

### Parameters

| Name                 | Type   | Required | Description              |
|----------------------|--------|----------|--------------------------|
| bk_biz_id            | int    | Yes      | Business ID              |
| name                 | string | Yes      | Cluster template name    |
| service_template_ids | array  | Yes      | Service template ID list |

### Request Example

```json
{
    "name": "test",
    "bk_biz_id": 20,
    "service_template_ids": [59]
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
        "id": 6,
        "name": "test",
        "bk_biz_id": 20,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-11-27T17:24:10.671658+08:00",
        "last_time": "2019-11-27T17:24:10.671658+08:00",
        "bk_supplier_account": "0"
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

#### data Field Explanation

| Name                | Type   | Description           |
|---------------------|--------|-----------------------|
| id                  | int    | Cluster template ID   |
| name                | array  | Cluster template name |
| bk_biz_id           | int    | Business ID           |
| creator             | string | Creator               |
| modifier            | string | Last modifier         |
| create_time         | string | Creation time         |
| last_time           | string | Update time           |
| bk_supplier_account | string | Supplier account      |
