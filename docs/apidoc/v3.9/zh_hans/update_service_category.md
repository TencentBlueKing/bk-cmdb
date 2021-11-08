### 功能描述

更新服务模板信息（目前仅名称字段可更新）

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| id            | int  | 否   | 服务模板ID |
| name         | string  | 否   | 服务模板名称 |


### 请求参数示例

```python
{
  "bk_biz_id": 1,
  "id": 3,
  "name": "222"
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "bk_biz_id": 1,
    "id": 3,
    "name": "222",
    "root_id": 3,
    "bk_supplier_account": "0",
    "is_built_in": false
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
| data | object | 更新后的服务分类信息 |

