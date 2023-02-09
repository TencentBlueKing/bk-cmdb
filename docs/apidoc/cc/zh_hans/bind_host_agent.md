### 功能描述

将agent绑定到主机上(v3.10.25+)

### 请求参数

{{ common_args_desc }}

### 请求参数

| 字段                |  类型              | 必选   |  描述                                   |
|---------------------|--------------------|--------|-----------------------------------------|
| list                  | array                | 是     | 要绑定的主机ID和agentID列表，最多200条    |

### list

| 字段                |  类型              | 必选   |  描述                                   |
|---------------------|--------------------|--------|-----------------------------------------|
| bk_host_id                  | int                | 是     | 要绑定agent的主机ID    |
| bk_agent_id                  | string                | 是     | 要绑定的agentID    |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "list": [
        {
            "bk_host_id": 1,
            "bk_agent_id": "xxxxxxxxxx"
        },
        {
            "bk_host_id": 2,
            "bk_agent_id": "yyyyyyyyyy"
        }
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807"
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