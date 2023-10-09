### 功能描述

查询主机对应的服务模版ID，该接口为节点管理专用接口，可能会随时调整，其它服务请勿使用(v3.10.11+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 参数       | 类型  | 必选 | 描述                |
| ---------- | ----- | ---- | ------------------- |
| bk_host_id | array | 是   | 主机id，最多为200个 |

#### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_id": [
        258,
        259
    ]
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
    "data": [
        {
            "bk_host_id": 258,
            "service_template_id": [
                3
            ]
        },
        {
            "bk_host_id": 259,
            "service_template_id": [
                1,
                2
            ]
        }
    ]
}
```

### 返回结果参数说明

#### response

| 字段                | 类型  | 描述       |
| ------------------- | ----- | ---------- |
| result     | bool   | 请求成功与否。true：请求成功；false：请求失败 |
| code       | int    | 错误编吗。0表示success，>0表示失败错误        |
| message    | string | 请求失败返回的错误信息                        |
| permission | object | 权限信息                                      |
| request_id | string | 请求链id                                      |
| data       | array  | 请求结果                                      |

#### data

| 字段                | 类型  | 描述       |
| ------------------- | ----- | ---------- |
| bk_host_id          | int   | 主机id     |
| service_template_id | array | 服务模版id |
