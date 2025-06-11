### 描述

删除动态分组 (版本：v3.9.6，权限：动态分组删除权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述       |
|-----------|--------|----|----------|
| bk_biz_id | int    | 是  | 业务ID     |
| id        | string | 是  | 动态分组主键ID |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |
