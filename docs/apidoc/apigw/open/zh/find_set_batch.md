### 描述

根据业务id和集群实例id列表，以及想要获取的属性列表，批量获取指定业务下集群的属性详情 (v3.8.6)

### 输入参数

| 参数名称      | 参数类型  | 必选 | 描述                              |
|-----------|-------|----|---------------------------------|
| bk_biz_id | int   | 是  | 业务ID                            |
| bk_ids    | array | 是  | 集群实例ID列表, 即bk_set_id列表，最多可填500个 |
| fields    | array | 是  | 集群属性列表，控制返回结果的集群信息里有哪些字段        |

### 调用示例

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        11,
        12
    ],
    "fields": [
        "bk_set_id",
        "bk_set_name",
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
            "bk_set_id": 12,
            "bk_set_name": "ss1",
            "create_time": "2020-05-15T22:15:51.67+08:00",
            "default": 0
        },
        {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "create_time": "2020-05-12T21:04:36.644+08:00",
            "default": 0
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

#### data

| 参数名称                 | 参数类型   | 描述                         |
|----------------------|--------|----------------------------|
| bk_set_name          | string | 集群名称                       |
| default              | int    | 0-普通集群，1-内置模块集合，默认为0       |
| bk_biz_id            | int    | 业务id                       |
| bk_capacity          | int    | 设计容量                       |
| bk_parent_id         | int    | 父节点的ID                     |
| bk_set_id            | int    | 集群id                       |
| bk_service_status    | string | 服务状态:1/2(1:开放,2:关闭)        |
| bk_set_desc          | string | 集群描述                       |
| bk_set_env           | string | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
| create_time          | string | 创建时间                       |
| last_time            | string | 更新时间                       |
| bk_supplier_account  | string | 开发商账号                      |
| description          | string | 数据的描述信息                    |
| set_template_version | array  | 集群模板的当前版本                  |
| set_template_id      | int    | 集群模板ID                     |
| bk_created_at        | string | 创建时间                       |
| bk_updated_at        | string | 更新时间                       |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
