### 描述

创建模块(权限：业务拓扑新建权限)

### 输入参数

| 参数名称            | 参数类型   | 必选 | 描述                                                       |
|-----------------|--------|----|----------------------------------------------------------|
| bk_biz_id       | int    | 是  | 业务ID                                                     |
| bk_set_id       | int    | 是  | 集群id                                                     |
| bk_parent_id    | int    | 是  | 父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于module一般指的是set的bk_set_id |
| bk_module_name  | string | 是  | 模块名                                                      |
| bk_module_type  | string | 否  | 模块类型                                                     |
| operator        | string | 否  | 主要维护人                                                    |
| bk_bak_operator | string | 否  | 备份维护人                                                    |
| bk_created_at   | string | 否  | 创建时间                                                     |
| bk_updated_at   | string | 否  | 更新时间                                                     |
| bk_created_by   | string | 否  | 创建人                                                      |
| bk_updated_by   | string | 否  | 更新人                                                      |

**注意：此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段，参数值的设置参考模块的属性字段配置
**

### 调用示例

```json
{
  "bk_parent_id": 4,
  "bk_module_name": "redis-1",
  "bk_module_type": "2",
  "operator": "admin",
  "bk_bak_operator": "admin",
  "bk_created_at": "",
  "bk_updated_at": "",
  "bk_created_by": "admin",
  "bk_updated_by": "admin"
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
    "bk_bak_operator": "admin",
    "bk_biz_id": 3,
    "bk_created_at": "2023-11-14T17:11:21.225+08:00",
    "bk_created_by": "admin",
    "bk_module_id": 20,
    "bk_module_name": "redis-1",
    "bk_module_type": "2",
    "bk_parent_id": 4,
    "bk_set_id": 4,
    "bk_supplier_account": "0",
    "bk_updated_at": "2023-11-14T17:11:21.225+08:00",
    "create_time": "2023-11-14T17:11:21.225+08:00",
    "default": 0,
    "host_apply_enabled": false,
    "last_time": "2023-11-14T17:11:21.225+08:00",
    "operator": "admin",
    "service_category_id": 2,
    "service_template_id": 0,
    "set_template_id": 0
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

| 参数名称                | 参数类型    | 描述           |
|---------------------|---------|--------------|
| bk_bak_operator     | string  | 备份维护人        |
| bk_module_id        | int     | 模型id         |
| bk_biz_id           | int     | 业务id         |
| bk_module_id        | int     | 主机所属的模块ID    |
| bk_module_name      | string  | 模块名          |
| bk_module_type      | string  | 模块类型         |
| bk_parent_id        | int     | 父节点的ID       |
| bk_set_id           | int     | 集群id         |
| bk_supplier_account | string  | 开发商账号        |
| create_time         | string  | 创建时间         |
| last_time           | string  | 更新时间         |
| default             | int     | 表示模块类型       |
| host_apply_enabled  | bool    | 是否启用主机属性自动应用 |
| operator            | string  | 主要维护人        |
| service_category_id | integer | 服务分类ID       |
| service_template_id | int     | 服务模版ID       |
| set_template_id     | int     | 集群模板ID       |
| bk_created_at       | string  | 创建时间         |
| bk_updated_at       | string  | 更新时间         |
| bk_created_by       | string  | 创建人          |
