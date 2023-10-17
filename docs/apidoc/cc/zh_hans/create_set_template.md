### 功能描述

在指定业务id下创建指定名称的集群模板，创建的集群模板通过指定的服务模板id去包含服务模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 | 类型   | 必选 | 描述           |
| -------------------- | ------ | ---- | -------------- |
| bk_biz_id            | int    | 是   | 业务ID         |
| name                 | string | 是   | 集群模板名称   |
| service_template_ids | array  | 是   | 服务模板ID列表 |


### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "name": "redis",
    "bk_biz_id": 3,
    "service_template_ids": [3]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "data": {
        "id": 1,
        "name": "redis",
        "bk_biz_id": 3,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2023-10-12T17:03:18.554918015+08:00",
        "last_time": "2023-10-12T17:03:18.554918015+08:00",
        "bk_supplier_account": "0"
    },
    "message": "success",
    "permission": null,
    "request_id": "a1841a8b69ba4c5ca2a65056133e6ffc"
}
```

### 返回结果参数说明

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                          |

#### data 字段说明

| 字段                | 类型   | 描述         |
| ------------------- | ------ | ------------ |
| id                  | int    | 集群模板ID   |
| name                | array  | 集群模板名称 |
| bk_biz_id           | int    | 业务ID       |
| creator             | string | 创建者       |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
