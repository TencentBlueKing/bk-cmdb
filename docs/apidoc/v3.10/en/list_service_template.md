### Functional description

list service templates

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| service_category_id         | int  | No   | Service Category ID |
| search         | string  | No   | search service templates by name，default is empty string |
| is_exact       | bool  | No   | whether exact match, used with search param, default is false (v3.9.19)|


### Request Parameters Example

```json
{
    "bk_biz_id": 1,
    "service_category_id": 1,
    "search": "test2",
    "is_exact": true,
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-name"
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "count": 1,
        "info": [
            {
                "bk_biz_id": 1,
                "id": 50,
                "name": "test2",
                "service_category_id": 1,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2019-09-18T20:31:29.607+08:00",
                "last_time": "2019-09-18T20:31:29.607+08:00",
                "bk_supplier_account": "0"
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
|id|integer|Service Template ID||
|name|array|Service Template Name||
|service_category_id|integer|Service Category ID||
