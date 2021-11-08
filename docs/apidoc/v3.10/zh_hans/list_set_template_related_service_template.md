### 功能描述

根据业务id和集群模板id,获取指定业务下某集群模版的服务模版列表

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选   |  描述                       |
|---------------------|------------|--------|-----------------------------|
| set_template_id     | int  | 是   | 集群模版ID |
| bk_biz_id           | int        | 是     | 业务ID |

### 请求参数示例

```json
{
  "set_template_id": 1,
  "bk_biz_id": 3
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
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
            "bk_supplier_account": "0"
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
            "bk_supplier_account": "0"
        }
    ]
}
```

### 返回结果参数说明

| 名称  | 类型  | 说明 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | object| 请求返回的数据 |

data 字段说明：

| 名称  | 类型  | 说明 |
|---|---|---|
| bk_biz_id           | int    | 业务ID       |
| id                  | int    | 服务模板ID   |
| name                | array  | 服务模板名称 |
| service_category_id | int    | 服务分类ID   |
| creator             | string | 创建者       |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
