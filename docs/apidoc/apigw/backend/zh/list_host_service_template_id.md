### 描述

查询主机对应的服务模版ID，该接口为节点管理专用接口，可能会随时调整，其它服务请勿使用(
版本：v3.10.11+，权限：主机池主机查看权限)

### 输入参数

| 参数名称       | 参数类型  | 必选 | 描述           |
|------------|-------|----|--------------|
| bk_host_id | array | 是  | 主机id，最多为200个 |

### 调用示例

```json
{
    "bk_host_id": [
        258,
        259
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
            "bk_host_id": 258,
            "service_template_id": [
                3
            ]
        },
        {
            "bk_host_id": 259,
            "service_template_id": [
                1,
                2
            ]
        }
    ]
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                          |
|------------|--------|-----------------------------|
| result     | bool   | 请求成功与否。true：请求成功；false：请求失败 |
| code       | int    | 错误编吗。0表示success，>0表示失败错误    |
| message    | string | 请求失败返回的错误信息                 |
| permission | object | 权限信息                        |
| data       | array  | 请求结果                        |

#### data

| 参数名称                | 参数类型  | 描述     |
|---------------------|-------|--------|
| bk_host_id          | int   | 主机id   |
| service_template_id | array | 服务模版id |
