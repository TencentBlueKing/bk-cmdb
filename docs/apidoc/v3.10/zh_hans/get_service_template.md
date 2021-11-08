### 功能描述

根据服务模板ID获取服务模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| service_template_id | int  | 是   | 服务模板ID |


### 请求参数示例

```json
{
  "bk_supplier_account": "0",
  "service_template_id": 51
}
```


### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "bk_biz_id": 3,
        "id": 51,
        "name": "mm2",
        "service_category_id": 12,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-05-26T09:46:15.259Z",
        "last_time": "2020-05-26T09:46:15.259Z",
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

| 字段|类型|说明|
|---|---|---|
|id|integer|服务模板ID|
|name|array|服务模板名称|
|service_category_id|integer|服务分类ID|
