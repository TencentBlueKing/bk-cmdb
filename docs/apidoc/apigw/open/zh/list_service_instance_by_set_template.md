### 描述

根据集群模版id获取服务实例列表

### 输入参数

| 参数名称            | 参数类型   | 必选 | 描述     |
|-----------------|--------|----|--------|
| bk_biz_id       | int    | 是  | 业务id   |
| set_template_id | int    | 是  | 集群模版ID |
| page            | object | 是  | 分页参数   |

#### page

| 参数名称  | 参数类型 | 必选 | 描述                     |
|-------|------|----|------------------------|
| start | int  | 是  | 记录开始位置                 |
| limit | int  | 是  | 每页限制条数,最大500,建议设定为200. |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "set_template_id":1,
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
                "bk_biz_id": 3,
                "id": 1,
                "name": "10.0.0.1_lgh-process-1",
                "labels": null,
                "service_template_id": 50,
                "bk_host_id": 1,
                "bk_module_id": 59,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2020-10-09T02:46:25.002Z",
                "last_time": "2020-10-09T02:46:25.002Z",
                "bk_supplier_account": "0"
            },
            {
                "bk_biz_id": 3,
                "id": 3,
                "name": "127.0.122.2_lgh-process-1",
                "labels": null,
                "service_template_id": 50,
                "bk_host_id": 3,
                "bk_module_id": 59,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2020-10-09T03:04:19.859Z",
                "last_time": "2020-10-09T03:04:19.859Z",
                "bk_supplier_account": "0"
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

#### data 字段说明

| 参数名称  | 参数类型  | 描述   |
|-------|-------|------|
| count | int   | 总数   |
| info  | array | 返回结果 |

#### info 字段说明

| 参数名称                | 参数类型   | 描述          |
|---------------------|--------|-------------|
| id                  | int    | 服务实例ID      |
| name                | string | 服务实例名称      |
| bk_biz_id           | int    | 业务id        |
| bk_module_id        | int    | 模型id        |
| service_template_id | int    | 服务模版ID      |
| labels              | map    | 标签信息        |
| bk_host_id          | int    | 主机id        |
| creator             | string | 本条数据创建者     |
| modifier            | string | 本条数据的最后修改人员 |
| create_time         | string | 创建时间        |
| last_time           | string | 更新时间        |
| bk_supplier_account | string | 开发商账号       |
