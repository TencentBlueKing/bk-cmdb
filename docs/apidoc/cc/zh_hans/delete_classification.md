### 功能描述

通过模型分类ID删除模型分类(权限：模型分组删除权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段 | 类型  | 必选 | 描述       |
|----|-----|----|----------|
| id | int | 是  | 分类数据记录ID |

**注意** 只能删除空模型分类，如果分类下有模型则删除失败

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id": 13
}
```

### 返回结果示例

#### 删除成功

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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
    "request_id": "8c6b89e7f0cb4fad836f55d50f81f2c6"
}
```

### 返回结果参数说明

#### response

| 字段         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |
