### Description

Update service category (Currently, only the name field can be updated. Permission: Service Category Editing Permission)

### Parameters

| Name      | Type   | Required | Description           |
|-----------|--------|----------|-----------------------|
| id        | int    | Yes      | Service category ID   |
| name      | string | Yes      | Service category name |
| bk_biz_id | int    | Yes      | Business ID           |

### Request Example

```python
{
  "bk_biz_id": 1,
  "id": 3,
  "name": "222"
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
        "bk_biz_id": 3,
        "id": 22,
        "name": "api",
        "bk_root_id": 21,
        "bk_parent_id": 21,
        "bk_supplier_account": "0",
        "is_built_in": false
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
| data       | object | Updated service category information                                |

#### data

| Name                | Type   | Description                      |
|---------------------|--------|----------------------------------|
| bk_biz_id           | int    | Business ID                      |
| id                  | int    | Service category ID              |
| name                | string | Service category name            |
| bk_root_id          | int    | Root service category ID         |
| bk_parent_id        | int    | Parent service category ID       |
| bk_supplier_account | string | Operator account                 |
| is_built_in         | bool   | Whether it is a built-in service |
