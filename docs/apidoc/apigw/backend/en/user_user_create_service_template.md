### Description

Creates a service template with the specified name and service class based on the provided service template name and
service class ID.

### Parameters

| Name                | Type   | Required | Description           |
|---------------------|--------|----------|-----------------------|
| name                | string | Yes      | Service template name |
| service_category_id | int    | Yes      | Service class ID      |
| bk_biz_id           | int    | Yes      | Business ID           |

### Request Example

```python
{
  "bk_biz_id": 1,
  "name": "test4",
  "service_category_id": 1
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
    "id": 52,
    "name": "test4",
    "service_category_id": 1,
    "creator": "admin",
    "modifier": "admin",
    "create_time": "2019-09-18T23:09:44.251970453+08:00",
    "last_time": "2019-09-18T23:09:44.251970568+08:00",
    "bk_supplier_account": "0"
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded; false: request failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                               |
| message    | string | Error message returned in case of request failure                                         |
| permission | object | Permission information                                                                    |
| data       | object | Data returned by the request                                                              |

#### Data Field Description

| Name                | Type   | Description                                  |
|---------------------|--------|----------------------------------------------|
| id                  | int    | Service template ID                          |
| bk_biz_id           | int    | Business ID                                  |
| name                | string | Service template name                        |
| service_category_id | int    | Service class ID                             |
| creator             | string | Creator of this data                         |
| modifier            | string | The last person to modify this piece of data |
| create_time         | string | Creation time                                |
| last_time           | string | Last modification time                       |
| bk_supplier_account | string | Developer account number                     |
