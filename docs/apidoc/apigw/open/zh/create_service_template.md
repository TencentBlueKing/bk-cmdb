### 描述

根据传入的服务模板名称和服务分类ID创建指定名称和服务分类的服务模板(权限：服务模板新建权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述     |
|---------------------|--------|----|--------|
| name                | string | 是  | 服务模板名称 |
| service_category_id | int    | 是  | 服务分类ID |
| bk_biz_id           | int    | 是  | 业务ID   |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "name": "test4",
  "service_category_id": 1
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
    "bk_biz_id": 1,
    "id": 52,
    "name": "test4",
    "service_category_id": 1,
    "creator": "admin",
    "modifier": "admin",
    "host_apply_enabled": false,
    "create_time": "2019-09-18T23:09:44.251970453+08:00",
    "last_time": "2019-09-18T23:09:44.251970568+08:00",
    "bk_supplier_account": "0"
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

| 参数名称                | 参数类型   | 描述           |
|---------------------|--------|--------------|
| id                  | int    | 服务模板ID       |
| bk_biz_id           | int    | 业务id         |
| name                | string | 服务模板名称       |
| service_category_id | int    | 服务模板ID       |
| creator             | string | 本条数据创建者      |
| modifier            | string | 本条数据的最后修改人员  |
| create_time         | string | 创建时间         |
| last_time           | string | 更新时间         |
| bk_supplier_account | string | 开发商账号        |
| host_apply_enabled  | bool   | 是否启用主机属性自动应用 |
