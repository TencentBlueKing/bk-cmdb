### 功能描述

根据传入的服务模板名称和服务分类ID创建指定名称和服务分类的服务模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| name            | string  | 是   | 服务模板名称 |
| service_category_id         | int  | 是   | 服务分类ID |
| bk_biz_id            | int  | 是   | 业务ID|

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "name": "redis-server",
    "service_category_id": 5
}
```

### 返回结果示例

```python
{
    "result": true,
    "code": 0,
    "data": {
        "bk_biz_id": 3,
        "id": 3,
        "name": "redis-server",
        "service_category_id": 5,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2023-10-12T16:58:59.667220357+08:00",
        "last_time": "2023-10-12T16:58:59.667220628+08:00",
        "bk_supplier_account": "0",
        "host_apply_enabled": false
    },
    "message": "success",
    "permission": null,
    "request_id": "f1416618038c4d11987f3883921cd0b3"
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
| data | object | 请求返回的数据 |

#### data 字段说明

| 字段|类型|描述|
|---|---|---|
|id|int|服务模板ID|
|bk_biz_id|int|业务id|
|name|string|服务模板名称|
|service_category_id|int|服务模板ID|
| creator              | string             | 本条数据创建者                                                                                 |
| modifier             | string             | 本条数据的最后修改人员            |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string       | 开发商账号 |
| host_apply_enabled|bool|是否启用主机属性自动应用|