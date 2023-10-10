### 功能描述

更新服务模板信息

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| name            | string  | 和service_category_id二选一必填，可都填   | 服务模板名称 |
| service_category_id            | int  | 和name二选一必填，可都填   | 服务分类id |
| id         | int  | 是   | 服务模板ID |
| bk_biz_id     |  int     | 是     | 业务ID |

### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "name": "test1",
  "id": 50,
  "service_category_id": 3
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
    "id": 50,
    "name": "test1",
    "service_category_id": 3,
    "creator": "admin",
    "modifier": "admin",
    "create_time": "2019-06-05T11:22:22.951+08:00",
    "last_time": "2019-06-05T11:22:22.951+08:00",
    "bk_supplier_account": "0",
    "host_apply_enabled": false
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
| data | object | 更新后的服务模板信息 |

#### data 字段说明

| 字段                | 类型   | 描述         |
| ------------------- | ------ | ------------ |
| id                  | int    | 服务模板ID   |
| name                | string  | 服务模板名称 |
| bk_biz_id           | int    | 业务ID       |
| service_category_id | int    | 服务分类id |
| creator             | string | 创建者       |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
| host_apply_enabled|bool|是否启用主机属性自动应用|
