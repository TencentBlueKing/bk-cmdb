### 描述

更新集群(权限：业务拓扑编辑权限)

### 输入参数

| 参数名称              | 参数类型   | 必选 | 描述                          |
|-------------------|--------|----|-----------------------------|
| bk_biz_id         | int    | 是  | 业务id                        |
| bk_set_id         | int    | 是  | 集群id                        |
| bk_set_name       | string | 否  | 集群名字                        |
| default           | int    | 否  | 0-普通集群，1-内置模块集合，默认为0        |
| set_template_id   | int    | 否  | 集群模板ID，需要通过集群模板创建集群时必填      |
| bk_capacity       | int    | 否  | 设计容量                        |
| description       | string | 否  | 备注、数据的描述信息                  |
| bk_set_desc       | string | 否  | 集群描述                        |
| bk_set_env        | string | 否  | 环境类型：测试(1)，体验(2)，正式(3, 默认值) |
| bk_service_status | string | 否  | 服务状态：开放(1, 默认值)，关闭(2)       |

**注意：此处仅对系统内置可编辑的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段；通过集群模板创建的集群，只能通过集群模板修改**

### 调用示例

```json
{
  "bk_set_name": "test",
  "default": 0,
  "bk_capacity": 500,
  "bk_set_desc": "集群描述",
  "description": "集群备注",
  "bk_set_env": "3",
  "bk_service_status": "1"
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
