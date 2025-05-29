### 描述

删除主机(权限:主机池主机删除权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述           |
|---------------------|--------|----|--------------|
| bk_supplier_account | string | 否  | 开发商账号        |
| bk_host_id          | string | 是  | 主机id，多个以逗号分隔 |

### 调用示例

```json
{
    "bk_host_id": "1,2,3"
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
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |
