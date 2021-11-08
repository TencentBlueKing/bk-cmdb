### Functional description

create service category

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| name            | string  | Yes   | service category name |
| parent_id         | int  | No   | parent node ID |


### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "name": "test101"
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
    "id": 6,
    "name": "test5",
    "root_id": 5,
    "parent_id": 5,
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
| data | object | new service category |

#### data description

| Field       | Type     | Description         |
|---|---|---|---|
|id|integer| service category ID||
|root_id|integer| root node ID||
|parent_id|integer| parent node ID||
|is_built_in|bool|is built in record or not|built in record is not editabled|

