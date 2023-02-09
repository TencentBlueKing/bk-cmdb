#### 功能描述

创建服务分类

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| name            | string  | 是   | 服务分类名称 |
| parent_id         | int  | 否   | 父节点ID |
| bk_biz_id         | int  | 是   | 业务ID |

### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "parent_id": 0,
  "bk_biz_id": 1,
  "name": "test101"
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "bk_biz_id": 1,
    "id": 6,
    "name": "test101",
    "root_id": 5,
    "parent_id": 5,
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
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | object | 新建的服务分类信息 |

#### data 字段说明

| 字段|类型|描述|
|---|---|---|
|id|integer|服务分类ID|
|root_id|integer|服务分类根节点ID|
|parent_id|integer|服务分类父节点ID|
|is_built_in|bool|是否是内置节点(内置节点不允许编辑)|
| bk_biz_id    | int     | 业务ID |
| name    | string     | 服务分类名称 |
| bk_supplier_account| string| 开发商账号|
