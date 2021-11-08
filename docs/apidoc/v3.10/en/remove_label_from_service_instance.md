### Functional description

remove label from service instance

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| instance_ids            | array  | Yes   | service instances ID array |
| keys            | array  | Yes   | key of lables to be remove |

### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "instance_ids": [60, 62],
  "keys": ["value1", "value3"]
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
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
