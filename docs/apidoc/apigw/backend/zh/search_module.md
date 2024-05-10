### 描述

查询模块

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述                  |
|---------------------|--------|----|---------------------|
| bk_supplier_account | string | 否  | 开发商账号               |
| bk_biz_id           | int    | 是  | 业务id                |
| bk_set_id           | int    | 否  | 集群ID                |
| fields              | array  | 是  | 查询字段，字段来自于模块定义的属性字段 |
| condition           | dict   | 是  | 查询条件，字段来自于模块定义的属性字段 |
| page                | dict   | 是  | 分页条件                |

#### page

| 参数名称  | 参数类型   | 必选 | 描述     |
|-------|--------|----|--------|
| start | int    | 是  | 记录开始位置 |
| limit | int    | 是  | 每页限制条数 |
| sort  | string | 否  | 排序字段   |

### 调用示例

```json
{
  "bk_biz_id": 2,
  "fields": [
    "bk_module_name",
    "bk_set_id"
  ],
  "condition": {
    "bk_module_name": "test"
  },
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
        "bk_module_name": "test",
        "bk_set_id": 11,
        "default": 0
      },
      {
        "bk_module_name": "test",
        "bk_set_id": 12,
        "default": 0
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
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data

| 参数名称  | 参数类型  | 描述                     |
|-------|-------|------------------------|
| count | int   | 数据数量                   |
| info  | array | 结果集，其中，所有字段均为模块定义的属性字段 |

#### info

| 参数名称                | 参数类型    | 描述           |
|---------------------|---------|--------------|
| bk_module_name      | string  | 模块名称         |
| bk_set_id           | int     | 集群id         |
| default             | int     | 表示模块类型       |
| bk_bak_operator     | string  | 备份维护人        |
| bk_module_id        | int     | 模型id         |
| bk_biz_id           | int     | 业务id         |
| bk_module_id        | int     | 主机所属的模块ID    |
| bk_module_type      | string  | 模块类型         |
| bk_parent_id        | int     | 父节点的ID       |
| bk_supplier_account | string  | 开发商账号        |
| create_time         | string  | 创建时间         |
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
