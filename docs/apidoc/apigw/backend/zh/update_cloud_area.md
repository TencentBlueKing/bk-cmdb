### 描述

更新管控区域(权限：管控区域编辑权限)

### 输入参数

| 参数名称          | 参数类型   | 必选 | 描述     |
|---------------|--------|----|--------|
| bk_cloud_id   | int    | 是  | 管控区域ID |
| bk_cloud_name | string | 否  | 管控区域名称 |

### 调用示例

```json
{
    "bk_cloud_id": 5,
    "bk_cloud_name": "管控区域1"
}

```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": "success"
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
