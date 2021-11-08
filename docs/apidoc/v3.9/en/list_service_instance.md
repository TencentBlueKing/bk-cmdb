### Functional description

list service instance list

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| bk_biz_id            | int  | Yes   | Biz ID |
| bk_module_id         | int  | No   | Module ID |
| selectors            | int  | No   | label filters，available operator values are: `=`,`!=`,`exists`,`!`,`in`,`notin`|
| page         | object  | No   | page parameter |
| search_key         | string  | No   | name filter |

### Request Parameters Example
```python

{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 1
  },
  "bk_module_id": 56,
  "search_key": "",
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  },{
    "key": "key1",
    "operator": "in",
    "values": ["value1", "value2"]
  }]
}

```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "count": 1,
    "info": [
      {
        "bk_biz_id": 1,
        "id": 72,
        "name": "t1",
        "bk_host_id": 26,
        "bk_module_id": 62,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-06-20T22:46:00.69+08:00",
        "last_time": "2019-06-20T22:46:00.69+08:00",
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
|id|integer|Service Instance ID||
|name|array|Service Instance Name||
|service_template_id|integer|Service Template ID||
|service_category_id|integer|Service Category ID||
|bk_host_id|integer|Host ID||
|bk_host_innerip|string|Host IP||
|bk_module_id|integer|Module ID||
|creator|string|Creator||
|modifier|string|Modifier||
|create_time|string|Create Time||
|last_time|string|Update Time||
|bk_supplier_account|string|Supplier Account ID||
