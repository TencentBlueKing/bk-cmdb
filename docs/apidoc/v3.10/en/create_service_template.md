### Functional description

create service template

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| name            | string  | Yes   | Service Template name |
| service_category_id         | int  | Yes   | Service Category ID |


### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "name": "test4",
  "service_category_id": 1
}
```

### Return Result Example

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

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |

#### data description

| Field       | Type     | Description         |
|---|---|---|---|
|count|integer|total count||
|info|array|response data||

#### info description

| Field       | Type     | Description         |
|---|---|---|---|
|id|integer|service template ID||
|name|array|service template||
|service_category_id|integer|service category ID||
