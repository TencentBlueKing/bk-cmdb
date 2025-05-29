### 描述

根据业务id和集群模板id,获取指定业务下某集群模版的服务模版列表

### 输入参数

| 参数名称            | 参数类型 | 必选 | 描述     |
|-----------------|------|----|--------|
| set_template_id | int  | 是  | 集群模版ID |
| bk_biz_id       | int  | 是  | 业务ID   |

### 调用示例

```json
{
  "set_template_id": 1,
  "bk_biz_id": 3
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
            "bk_biz_id": 3,
            "id": 48,
            "name": "sm1",
            "service_category_id": 2,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:14:57.691Z",
            "last_time": "2020-05-15T14:14:57.691Z",
            "bk_supplier_account": "0",
            "host_apply_enabled": false
        },
        {
            "bk_biz_id": 3,
            "id": 49,
            "name": "sm2",
            "": 16,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:19:09.813Z",
            "last_time": "2020-05-15T14:19:09.813Z",
            "bk_supplier_account": "0",
            "host_apply_enabled": false
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

data 字段说明：

| 参数名称                | 参数类型   | 描述           |
|---------------------|--------|--------------|
| bk_biz_id           | int    | 业务ID         |
| id                  | int    | 服务模板ID       |
| name                | string | 服务模板名称       |
| service_category_id | int    | 服务分类ID       |
| creator             | string | 创建者          |
| modifier            | string | 最后修改人员       |
| create_time         | string | 创建时间         |
| last_time           | string | 更新时间         |
| bk_supplier_account | string | 开发商账号        |
| host_apply_enabled  | bool   | 是否启用主机属性自动应用 |
