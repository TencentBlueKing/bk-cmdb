### Functional description

list service caegory

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |


### Request Parameters Example

```python
{
  "bk_biz_id": 1,
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
    "count": 20,
    "info": [
      {
	"bk_biz_id": 1,
        "id": 16,
        "name": "Apache",
        "bk_root_id": 14,
        "bk_parent_id": 14,
        "bk_supplier_account": "0",
        "is_built_in": true
      },
      {
	"bk_biz_id": 1,
        "id": 19,
        "name": "Ceph",
        "bk_root_id": 18,
        "bk_parent_id": 18,
        "bk_supplier_account": "0",
        "is_built_in": true
      },
      {
	"bk_biz_id": 1,
        "id": 1,
        "name": "Default",
        "bk_root_id": 1,
        "bk_supplier_account": "0",
        "is_built_in": true
      }
    ]
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

#### Data field description

| Field       | Type     | Description         |
|---|---|---|---|
|count|integer|total count||
|info|array|response data||

#### Info field description

| Field       | Type     | Description         |
|---|---|---|---|
|id|integer|service category ID||
|name|string|service category name||
|bk_root_id|integer|root node ID||
|bk_parent_id|integer|parent node ID||
|is_built_in|bool|is built in||
