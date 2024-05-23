### 描述

创建服务分类(权限：服务分类新建权限)

### 输入参数

| 参数名称         | 参数类型   | 必选 | 描述     |
|--------------|--------|----|--------|
| name         | string | 是  | 服务分类名称 |
| bk_parent_id | int    | 否  | 父节点ID  |
| bk_biz_id    | int    | 是  | 业务ID   |

### 调用示例

```json
{
  "bk_parent_id": 0,
  "bk_biz_id": 1,
  "name": "test101"
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
    "bk_biz_id": 1,
    "id": 6,
    "name": "test101",
    "bk_root_id": 5,
    "bk_parent_id": 5,
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
| data       | object | 新建的服务分类信息                  |

#### data 字段说明

| 参数名称                | 参数类型    | 描述                 |
|---------------------|---------|--------------------|
| id                  | integer | 服务分类ID             |
| root_id             | integer | 服务分类根节点ID          |
| parent_id           | integer | 服务分类父节点ID          |
| is_built_in         | bool    | 是否是内置节点(内置节点不允许编辑) |
| bk_biz_id           | int     | 业务ID               |
| name                | string  | 服务分类名称             |
| bk_supplier_account | string  | 开发商账号              |
