### 描述

更新模块(权限：业务拓扑编辑权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述    |
|---------------------|--------|----|-------|
| bk_supplier_account | string | 否  | 开发商账号 |
| bk_biz_id           | int    | 是  | 业务id  |
| bk_set_id           | int    | 是  | 集群id  |
| bk_module_id        | int    | 是  | 模块id  |
| data                | dict   | 是  | 模块数据  |

#### data

| 参数名称            | 参数类型   | 必选 | 描述    |
|-----------------|--------|----|-------|
| bk_module_name  | string | 否  | 模块名   |
| bk_module_type  | string | 否  | 模块类型  |
| operator        | string | 否  | 主要维护人 |
| bk_bak_operator | string | 否  | 备份维护人 |

**注意：此处data参数仅对系统内置可编辑的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段；通过服务模板创建的模块，只能通过服务模板修改
**

### 调用示例

```json
{
    "bk_biz_id": 1,
    "bk_set_id": 1,
    "bk_module_id": 1,
    "data": {
        "bk_module_name": "test",
        "bk_module_type": "1",
        "operator": "admin",
        "bk_bak_operator": "admin"
    }
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
