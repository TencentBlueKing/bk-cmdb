### Functional description

update service template info

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| id            | int  | No   | Service Template ID |
| name         | string  | No   | Service Template Name |

### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "id": 3,
  "name": "222"
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "bk_biz_id": 1,
    "id": 3,
    "name": "222",
    "root_id": 3,
    "bk_supplier_account": "0",
    "is_built_in": false
  }
}
```

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |
