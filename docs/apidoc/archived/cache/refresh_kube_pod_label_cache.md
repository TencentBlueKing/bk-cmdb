### 协议

POST /api/v3/cache/refresh/kube/pod/label

### 描述

刷新业务中Pod的标签键值缓存 (版本：v3.13.5+，权限：业务访问)

**注意：**
- 该接口为异步接口，调用时会直接返回，通过后台任务刷新缓存数据。

### 输入参数

| 参数名称      | 参数类型 | 必选 | 描述   |
|-----------|------|----|------|
| bk_biz_id | int  | 是  | 业务ID |

### 调用示例

```json
{
  "bk_biz_id": 3
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
