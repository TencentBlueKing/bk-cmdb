### 描述

根据业务ID和模块实例ID列表，加上想要返回的模块属性列表，批量获取指定业务下模块实例的属性信息 (v3.8.6)

### 输入参数

| 参数名称      | 参数类型  | 必选 | 描述                                 |
|-----------|-------|----|------------------------------------|
| bk_biz_id | int   | 是  | 业务ID                               |
| bk_ids    | array | 是  | 模块实例ID列表, 即bk_module_id列表，最多可填500个 |
| fields    | array | 是  | 模块属性列表，控制返回结果的模块信息里有哪些字段           |

### 调用示例

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        56,
        57,
        58,
        59,
        60
    ],
    "fields": [
        "bk_module_id",
        "bk_module_name",
        "create_time"
    ]
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [
        {
            "bk_module_id": 60,
            "bk_module_name": "sm1",
            "create_time": "2020-05-15T22:15:51.725+08:00",
            "default": 0
        },
        {
            "bk_module_id": 59,
            "bk_module_name": "m1",
            "create_time": "2020-05-12T21:04:47.286+08:00",
            "default": 0
        },
        {
            "bk_module_id": 58,
            "bk_module_name": "待回收",
            "create_time": "2020-05-12T21:03:37.238+08:00",
            "default": 3
        },
        {
            "bk_module_id": 57,
            "bk_module_name": "故障机",
            "create_time": "2020-05-12T21:03:37.183+08:00",
            "default": 2
        },
        {
            "bk_module_id": 56,
            "bk_module_name": "空闲机",
            "create_time": "2020-05-12T21:03:37.122+08:00",
            "default": 1
        }
    ]
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | array  | 请求返回的数据                    |

#### data说明

| 参数名称                | 参数类型    | 描述           |
|---------------------|---------|--------------|
| bk_module_id        | int     | 模块id         |
| bk_module_name      | string  | 模块名称         |
| default             | int     | 表示模块类型       |
| create_time         | string  | 创建时间         |
| bk_set_id           | int     | 集群id         |
| bk_bak_operator     | string  | 备份维护人        |
| bk_biz_id           | int     | 业务id         |
| bk_module_type      | string  | 模块类型         |
| bk_parent_id        | int     | 父节点的ID       |
| bk_supplier_account | string  | 开发商账号        |
| last_time           | string  | 更新时间         |
| host_apply_enabled  | bool    | 是否启用主机属性自动应用 |
| operator            | string  | 主要维护人        |
| service_category_id | integer | 服务分类ID       |
| service_template_id | int     | 服务模版ID       |
| set_template_id     | int     | 集群模板ID       |
| bk_created_at       | string  | 创建时间         |
| bk_updated_at       | string  | 更新时间         |
| bk_created_by       | string  | 创建人          |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
