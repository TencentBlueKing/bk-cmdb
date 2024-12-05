### 请求方式

GET /migrate/v3/find/system/tenant_db_relation

### 描述

查询租户与DB的关系

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": [
    {
      "tenant_id": "blueking",
      "database": "masteruuid"
    }
  ]
}
```

### 响应参数说明

| 参数名称       | 参数类型         | 描述                         |
|------------|--------------|----------------------------|
| result     | bool         | 请求成功与否。true:请求成功；false请求失败 |
| code       | int          | 错误编码。 0表示success，>0表示失败错误  |
| message    | string       | 请求失败返回的错误信息                |
| permission | object       | 权限信息                       |
| data       | object array | 请求返回的数据                    |

#### data[n]

| 参数名称      | 参数类型   | 描述      |
|-----------|--------|---------|
| tenant_id | string | 租户ID    |
| database  | string | 数据库唯一标识 |
