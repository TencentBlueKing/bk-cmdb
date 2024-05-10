### 描述

更新项目id，此接口为BCS进行项目数据迁移时的专用接口，其他平台请勿使用(版本：v3.10.23+，权限：项目更新权限)

### 输入参数

| 参数名称          | 参数类型   | 必选 | 描述                    |
|---------------|--------|----|-----------------------|
| id            | int    | 是  | project在cc中的id唯一标识    |
| bk_project_id | string | 是  | bk_project_id需要更新的最终值 |

### 调用示例

```json
{
    "id": 1,
    "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
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
