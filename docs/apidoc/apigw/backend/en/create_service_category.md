### Description

Create Service Category (Permission: Service Category Creation Permission)

### Parameters

| Name         | Type   | Required | Description           |
|--------------|--------|----------|-----------------------|
| name         | string | Yes      | Service category name |
| bk_parent_id | int    | No       | Parent node ID        |
| bk_biz_id    | int    | Yes      | Business ID           |

### Request Example

```python
{
  "bk_parent_id": 0,
  "bk_biz_id": 1,
  "name": "test101"
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
    "id": 6,
    "name": "test101",
    "bk_root_id": 5,
    "bk_parent_id": 5,
    "bk_supplier_account": "0",
    "is_built_in": false
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Newly created service category information                                  |

#### data

| Name                | Type    | Description                                                     |
|---------------------|---------|-----------------------------------------------------------------|
| id                  | integer | Service category ID                                             |
| root_id             | integer | Service category root node ID                                   |
| parent_id           | integer | Service category parent node ID                                 |
| is_built_in         | bool    | Whether it is a built-in node (built-in nodes cannot be edited) |
| bk_biz_id           | int     | Business ID                                                     |
| name                | string  | Service category name                                           |
| bk_supplier_account | string  | Vendor account                                                  |
