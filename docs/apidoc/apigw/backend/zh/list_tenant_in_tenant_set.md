### 描述

查询租户集中的租户(版本：v3.15.1+，权限：租户集访问权限)

### 输入参数

| 参数名称          | 参数类型 | 必选 | 描述    |
|---------------|------|----|-------|
| tenant_set_id | int  | 是  | 租户集id |

### 调用示例

```json
{
  "tenant_set_id": 1
}
```

### 响应示例

```json
{
  "data": [
    {
      "id": "default",
      "name": "Default",
      "status": "enabled"
    },
    {
      "id": "test",
      "name": "Test",
      "status": "disabled"
    }
  ]
}
```

### 响应参数说明

| 参数名称   | 参数类型   | 描述                              |
|--------|--------|---------------------------------|
| id     | string | 租户 ID                           |
| name   | string | 租户名                             |
| status | string | 租户状态，enabled 表示启用，disabled 表示禁用 |