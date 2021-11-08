### Functional description

batch delete service instances

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| service_instance_ids | array  | Yes   | Service Instance IDs |

### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "service_instance_ids": [48]
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
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
