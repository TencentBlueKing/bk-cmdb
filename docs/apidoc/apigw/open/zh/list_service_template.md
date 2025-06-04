### 描述

根据业务id查询服务模板列表,可加上服务分类id进一步查询

### 输入参数

| 参数名称                 | 参数类型      | 必选 | 描述                                                         |
|----------------------|-----------|----|------------------------------------------------------------|
| bk_biz_id            | int       | 是  | 业务ID                                                       |
| service_category_id  | int       | 否  | 服务分类ID                                                     |
| search               | string    | 否  | 按服务模版名查询，默认为空                                              |
| is_exact             | bool      | 否  | 是否精确匹配服务模版名，默认为否，和search参数搭配使用，在search参数不为空的情况下有效（v3.9.19） |
| service_template_ids | int array | 否  | 服务模板id                                                     |
| page                 | object    | 是  | 分页参数                                                       |

#### page

| 参数名称  | 参数类型   | 必选 | 描述           |
|-------|--------|----|--------------|
| start | int    | 是  | 记录开始位置       |
| limit | int    | 是  | 每页限制条数,最大500 |
| sort  | string | 否  | 排序字段         |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "service_category_id": 1,
    "service_template_ids":[5,6],
    "search": "test2",
    "is_exact": true,
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-name"
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
        "count": 1,
        "info": [
            {
                "bk_biz_id": 1,
                "id": 50,
                "name": "test2",
                "service_category_id": 1,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2019-09-18T20:31:29.607+08:00",
                "last_time": "2019-09-18T20:31:29.607+08:00",
                "bk_supplier_account": "0",
                "host_apply_enabled": false
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

| 参数名称                | 参数类型    | 描述           |
|---------------------|---------|--------------|
| bk_biz_id           | int     | 业务id         |
| id                  | int     | 服务模板ID       |
| name                | array   | 服务模板名称       |
| service_category_id | integer | 服务分类ID       |
| creator             | string  | 创建人          |
| modifier            | string  | 修改人          |
| create_time         | string  | 创建时间         |
| last_time           | string  | 修复时间         |
| bk_supplier_account | string  | 供应商ID        |
| host_apply_enabled  | bool    | 是否启用主机属性自动应用 |
