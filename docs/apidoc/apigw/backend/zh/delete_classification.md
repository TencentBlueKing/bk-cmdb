### 描述

通过模型分类ID删除模型分类(权限：模型分组删除权限)

### 输入参数

| 参数名称 | 参数类型 | 必选 | 描述       |
|------|------|----|----------|
| id   | int  | 是  | 分类数据记录ID |

**注意** 只能删除空模型分类，如果分类下有模型则删除失败

### 调用示例

```json
{
    "id": 13
}
```

### 响应示例

#### 删除成功

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": "success"
}
```

#### 分类下有模型，删除失败

```json
{
    "result": false,
    "code": 1101029,
    "data": null,
    "message": "There is a model under the category, not allowed to delete",
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
