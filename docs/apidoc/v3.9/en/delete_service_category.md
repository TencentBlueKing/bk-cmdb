### Functional description

delete service category

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| id            | int  | Yes   | service category ID |

### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "id": 6
}
```

### Return Result Example

```python
{
  "result": false,
  "code": 1199054,
  "message": "operate built-in object forbidden",
  "data": ""
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
