### 功能描述

根据业务id查询集群模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                | 类型   | 必选 | 描述           |
| ------------------- | ------ | ---- | -------------- |
| bk_biz_id           | int    | 是   | 业务ID         |
| set_template_ids    | array  | 否   | 集群模板ID数组 |
| page                | object | 否   | 分页信息       |

#### page 字段说明

| 字段  | 类型   | 必选 | 描述                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | 否   | 记录开始位置          |
| limit | int    | 否   | 每页限制条数,最大1000 |
| sort  | string | 否   | 排序字段，'-'表示倒序 |


### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "bk_biz_id": 3,
    "set_template_ids": [
        3
    ],
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-name"
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "data": {
        "count": 1,
        "info": [
            {
                "id": 3,
                "name": "redis",
                "bk_biz_id": 3,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2023-10-12T10:14:40.71Z",
                "last_time": "2023-10-12T10:14:40.71Z",
                "bk_supplier_account": "0"
            }
        ]
    },
    "message": "success",
    "permission": null,
    "request_id": "2fb8761b27394865ba79dd7b7de09636"
}
```

### 返回结果参数说明

#### response

| 名称    | 类型   | 描述                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

#### data 字段说明

| 字段  | 类型  | 说明     |
| ----- | ----- | -------- |
| count | int   | 总数     |
| info  | array | 返回结果 |

#### info 字段说明

| 字段                | 类型   | 描述         |
| ------------------- | ------ | ------------ |
| id                  | int    | 集群模板ID   |
| name                | array  | 集群模板名称 |
| bk_biz_id           | int    | 业务ID       |
| creator             | string | 创建者       |
| modifier            | string | 最后修改人员 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string | 开发商账号   |
