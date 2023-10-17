### 功能描述

获取动态分组详情 (V3.9.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | 是     | 业务ID |
| id        |  string  | 是     | 目标动态分组主键ID |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "data": {
        "bk_biz_id": 3,
        "id": "293dcda1-68da-11ee-b1d6-52540075cc6d",
        "name": "主机名含：host",
        "bk_obj_id": "host",
        "info": {
            "condition": [
                {
                    "bk_obj_id": "set",
                    "condition": [
                        {
                            "field": "bk_set_name",
                            "operator": "$in",
                            "value": [
                                "aaTset"
                            ]
                        }
                    ]
                },
                {
                    "bk_obj_id": "module",
                    "condition": [
                        {
                            "field": "bk_module_name",
                            "operator": "$in",
                            "value": [
                                "aaMod"
                            ]
                        }
                    ]
                },
                {
                    "bk_obj_id": "host",
                    "condition": [
                        {
                            "field": "bk_host_name",
                            "operator": "$regex",
                            "value": "host"
                        }
                    ]
                }
            ]
        },
        "create_user": "admin",
        "modify_user": "admin",
        "create_time": "2023-10-12T08:34:29.913Z",
        "last_time": "2023-10-12T08:34:29.913Z"
    },
    "message": "success",
    "permission": null,
    "request_id": "0bedb6e7ca594279a3670ddc2a80ce5c"
}
```

### 返回结果参数
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| bk_biz_id    | int     | 业务ID |
| id           | string  | 动态分组主键ID |
| bk_obj_id    | string  | 动态分组的目标资源对象类型,目前可以为host,set |
| name         | string  | 动态分组命名 |
| info         | object  | 动态分组规则信息 |
| last_time    | string  | 更新时间 |
| modify_user  | string  | 修改者 |
| create_time  | string  | 创建时间 |
| create_user  | string  | 创建者 |

#### data.info.condition

| 字段      |  类型     |  描述      |
|-----------|-----------|------------|
| bk_obj_id |  string   | 条件对象资源类型, host类型的动态分组支持的info.conditon:set,module,host；set类型的动态分组支持的info.condition:set |
| condition |  array    | 查询条件 |

#### data.info.condition.condition

| 字段      |  类型     |  描述       |
|-----------|------------|------------|
| field     |  string    | 对象的字段 |
| operator  |  string    | 操作符, op值为eq(相等)/ne(不等)/in(属于)/nin(不属于) |
| value     |  object    | 字段对应的值 |
