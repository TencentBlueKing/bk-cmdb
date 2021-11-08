### 功能描述

在指定业务id下创建指定名称的集群模板，创建的集群模板通过指定的服务模板id去包含服务模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 | 类型   | 必选 | 描述           |
| -------------------- | ------ | ---- | -------------- |
| bk_biz_id            | int    | 是   | 业务ID         |
| name                 | string | 是   | 集群模板名称   |
| service_template_ids | array  | 是   | 服务模板ID列表 |


### 请求参数示例

```json
{
    "bk_supplier_account": "0",
    "name": "test",
    "bk_biz_id": 20,
    "service_template_ids": [59]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "id": 6,
        "name": "test",
        "bk_biz_id": 20,
        "version": 0,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-11-27T17:24:10.671658+08:00",
        "last_time": "2019-11-27T17:24:10.671658+08:00",
        "bk_supplier_account": "0"
    }
}
```

### 返回结果参数说明

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| data    | object | 请求返回的数据                          |

#### data 字段说明

| 字段                | 类型   | 描述         |
| ------------------- | ------ | ------------ |
| id                  | int    | 集群模板ID   |
| name                | array  | 集群模板名称 |
| bk_biz_id           | int    | 业务ID       |
| version             | int    | 集群模板版本 |
| creator             | string | 创建者       |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
