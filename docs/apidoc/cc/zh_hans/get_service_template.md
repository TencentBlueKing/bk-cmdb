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
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "service_template_id": 51
}
```


### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "bk_biz_id": 3,
        "id": 51,
        "name": "mm2",
        "service_category_id": 12,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-05-26T09:46:15.259Z",
        "last_time": "2020-05-26T09:46:15.259Z",
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
| data | object | 请求返回的数据 |

#### data 字段说明

| 字段|类型|说明|
|---|---|---|
|bk_biz_id|int|业务ID|
|id|int|服务模板ID|
|name|array|服务模板名称|
|service_category_id|integer|服务分类ID|
| creator             | string | 创建者       |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
| host_apply_enabled|bool|是否启用主机属性自动应用|
