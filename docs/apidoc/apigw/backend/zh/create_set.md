### 描述

创建集群(权限：业务拓扑新建权限)

### 输入参数

| 参数名称              | 参数类型   | 必选 | 描述                                           |
|-------------------|--------|----|----------------------------------------------|
| bk_biz_id         | int    | 是  | 业务ID                                         |
| bk_parent_id      | int    | 是  | 父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于set一般指的是业务ID |
| bk_set_name       | string | 是  | 集群名字                                         |
| default           | int    | 否  | 0-普通集群，1-内置模块集合，默认为0                         |
| set_template_id   | int    | 否  | 集群模板ID，需要通过集群模板创建集群时必填                       |
| bk_capacity       | int    | 否  | 设计容量                                         |
| description       | string | 否  | 备注、数据的描述信息                                   |
| bk_set_desc       | string | 否  | 集群描述                                         |
| bk_set_env        | string | 否  | 环境类型：测试(1)，体验(2)，正式(3, 默认值)                  |
| bk_service_status | string | 否  | 服务状态：开放(1, 默认值)，关闭(2)                        |
| bk_created_at     | string | 否  | 创建时间                                         |
| bk_updated_at     | string | 否  | 更新时间                                         |
| bk_created_by     | string | 否  | 创建人                                          |
| bk_updated_by     | string | 否  | 更新人                                          |

**注意：此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段，参数值的设置参考集群的属性字段配置
**

### 调用示例

```json
{
  "bk_parent_id": 3,
  "bk_set_name": "set_a1",
  "set_template_id": 0,
  "default": 0,
  "bk_capacity": 1000,
  "bk_set_desc": "test-set",
  "bk_set_env": "1",
  "bk_service_status": "1",
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
    "bk_biz_id": 3,
    "bk_capacity": 1000,
    "bk_created_at": "2023-11-14T17:30:43.048+08:00",
    "bk_created_by": "admin",
    "bk_parent_id": 3,
    "bk_service_status": "1",
    "bk_set_desc": "test-set",
    "bk_set_env": "1",
    "bk_set_id": 10,
    "bk_set_name": "set_a1",
    "bk_supplier_account": "0",
    "bk_updated_at": "2023-11-14T17:30:43.048+08:00",
    "create_time": "2023-11-14T17:30:43.048+08:00",
    "default": 0,
    "description": "",
    "last_time": "2023-11-14T17:30:43.048+08:00",
    "set_template_id": 0,
    "set_template_version": null
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

| 参数名称                 | 参数类型   | 描述                         |
|----------------------|--------|----------------------------|
| bk_biz_id            | int    | 业务id                       |
| bk_capacity          | int    | 设计容量                       |
| bk_parent_id         | int    | 父节点的ID                     |
| bk_set_id            | int    | 集群id                       |
| bk_service_status    | string | 服务状态:1/2(1:开放,2:关闭)        |
| bk_set_desc          | string | 集群描述                       |
| bk_set_env           | string | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
| bk_set_name          | string | 集群名称                       |
| create_time          | string | 创建时间                       |
| last_time            | string | 更新时间                       |
| bk_supplier_account  | string | 开发商账号                      |
| default              | int    | 0-普通集群，1-内置模块集合，默认为0       |
| description          | string | 数据的描述信息                    |
| set_template_version | array  | 集群模板的当前版本                  |
| set_template_id      | int    | 集群模板ID                     |
| bk_created_at        | string | 创建时间                       |
| bk_updated_at        | string | 更新时间                       |
| bk_created_by        | string | 创建人                        |
