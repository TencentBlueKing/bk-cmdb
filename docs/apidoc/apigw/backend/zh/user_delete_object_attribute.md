### 描述

删除对象模型属性，可以删除业务自定义字段(权限：模型编辑权限)

### 输入参数

| 参数名称 | 参数类型 | 必选 | 描述              |
|------|------|----|-----------------|
| id   | int  | 是  | 被删除的数据记录的唯一标识ID |

### 调用示例

```json
{
    "id" : 0
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
