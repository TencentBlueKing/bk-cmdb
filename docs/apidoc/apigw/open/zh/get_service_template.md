### 描述

根据服务模板ID获取服务模板

### 输入参数

| 参数名称                | 参数类型 | 必选 | 描述     |
|---------------------|------|----|--------|
| service_template_id | int  | 是  | 服务模板ID |

### 调用示例

```json
{
  "service_template_id": 51
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
        "id": 51,
        "name": "mm2",
        "service_category_id": 12,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-05-26T09:46:15.259Z",
        "last_time": "2020-05-26T09:46:15.259Z",
        "bk_supplier_account": "0",
        "host_apply_enabled": false
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

| 参数名称                | 参数类型    | 描述           |
|---------------------|---------|--------------|
| bk_biz_id           | int     | 业务ID         |
| id                  | int     | 服务模板ID       |
| name                | array   | 服务模板名称       |
| service_category_id | integer | 服务分类ID       |
| creator             | string  | 创建者          |
| modifier            | string  | 最后修改人员       |
| create_time         | string  | 创建时间         |
| last_time           | string  | 更新时间         |
| bk_supplier_account | string  | 开发商账号        |
| host_apply_enabled  | bool    | 是否启用主机属性自动应用 |
