### 描述

将agent绑定到主机上(版本：v3.10.25+，权限：主机AgentID管理权限)

### 输入参数

| 参数名称 | 参数类型  | 必选 | 描述                        |
|------|-------|----|---------------------------|
| list | array | 是  | 要绑定的主机ID和agentID列表，最多200条 |

#### list

| 参数名称        | 参数类型   | 必选 | 描述            |
|-------------|--------|----|---------------|
| bk_host_id  | int    | 是  | 要绑定agent的主机ID |
| bk_agent_id | string | 是  | 要绑定的agentID   |

### 调用示例

```json
{
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

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
