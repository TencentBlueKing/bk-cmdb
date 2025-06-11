### 描述

根据条件查询业务下的模块 (v3.9.7)

### 输入参数

| 参数名称                    | 参数类型   | 必选 | 描述                       |
|-------------------------|--------|----|--------------------------|
| bk_biz_id               | int    | 是  | 业务ID                     |
| bk_set_ids              | array  | 否  | 集群ID列表, 最多可填200个         |
| bk_service_template_ids | array  | 否  | 服务模板ID列表                 |
| fields                  | array  | 是  | 模块属性列表，控制返回结果的模块信息里有哪些字段 |
| page                    | object | 是  | 分页信息                     |

#### page 字段说明

| 参数名称  | 参数类型 | 必选 | 描述           |
|-------|------|----|--------------|
| start | int  | 是  | 记录开始位置       |
| limit | int  | 是  | 每页限制条数,最大500 |

### 调用示例

```json
{
    "bk_biz_id": 2,
    "bk_set_ids":[1,2],
    "bk_service_template_ids": [3,4],
    "fields":["bk_module_id", "bk_module_name"],
    "page": {
        "start": 0,
        "limit": 10
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
    "data": {
        "count": 2,
        "info": [
            {
                "bk_module_id": 8,
                "bk_module_name": "license"
            },
            {
                "bk_module_id": 12,
                "bk_module_name": "gse_proc"
            }
        ]
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
| data       | object | 请求返回的数据                    |

#### data 字段说明：

| 参数名称  | 参数类型         | 描述     |
|-------|--------------|--------|
| count | int          | 记录条数   |
| info  | object array | 模块实际数据 |

#### data.info 字段说明:

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
