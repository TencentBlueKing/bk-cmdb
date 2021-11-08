### Functional description

add label for service instance

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required    | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
|instance_ids|array|Yes|Service Instance ID|
|labels|object|Yes|Labels to be add|

#### labels field description
- key field validation rule: `^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- value field validation rule: `^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "instance_ids": [59, 62],
  "labels": {
    "key1": "value1",
    "key2": "value2"
  }
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
