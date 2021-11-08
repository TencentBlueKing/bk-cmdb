### Functional description

delete service template

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| service_template_id | int  | Yes   | Service Template ID |

### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "service_template_id": 1
}
```

### Return Result Example

```python
{
  "result": false,
  "code": 1199056,
  "message": "remove referenced record forbidden",
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
