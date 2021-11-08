### 功能描述

根据传入的服务模板名称的服务分类ID创建指定名称和服务分类的服务模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| name            | string  | 是   | 服务模板名称 |
| service_category_id         | int  | 是   | 服务分类ID |


### 请求参数示例

```python
{
  "bk_biz_id": 1,
  "name": "test4",
  "service_category_id": 1
}
```

### 返回结果示例

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

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | object | 请求返回的数据 |

#### data 字段说明

| 字段|类型|说明|描述|
|---|---|---|---|
|count|integer|总数||
|info|array|返回结果||

#### info 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|integer|服务模板ID||
|name|array|服务模板名称||
|service_category_id|integer|服务模板ID||
