### 描述

更新服务分类(目前仅名称字段可更新，权限：服务分类编辑权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述     |
|-----------|--------|----|--------|
| id        | int    | 是  | 服务分类ID |
| name      | string | 是  | 服务分类名称 |
| bk_biz_id | int    | 是  | 业务ID   |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "id": 3,
  "name": "222"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "bk_biz_id": 3,
        "id": 22,
        "name": "api",
        "bk_root_id": 21,
        "bk_parent_id": 21,
        "bk_supplier_account": "0",
        "is_built_in": false
    }
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 更新后的服务分类信息                 |

#### data

| 参数名称                | 参数类型   | 描述      |
|---------------------|--------|---------|
| bk_biz_id           | int    | 业务id    |
| id                  | int    | 服务分类id  |
| name                | string | 服务分类名称  |
| bk_root_id          | int    | 根服务分类id |
| bk_parent_id        | int    | 父服务分类id |
| bk_supplier_account | string | 运营商账号   |
| is_built_in         | bool   | 是否是内置服务 |
