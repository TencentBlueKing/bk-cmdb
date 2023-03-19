### 功能描述

根据业务ID和模块实例ID列表，加上想要返回的模块属性列表，批量获取指定业务下模块实例的属性信息 (v3.8.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  | 是     | 业务ID |
| bk_ids  |  array  | 是     | 模块实例ID列表, 即bk_module_id列表，最多可填500个 |
| fields  |   array   | 是     | 模块属性列表，控制返回结果的模块信息里有哪些字段 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "bk_ids": [
        56,
        57,
        58,
        59,
        60
    ],
    "fields": [
        "bk_module_id",
        "bk_module_name",
        "create_time"
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_module_id": 60,
            "bk_module_name": "sm1",
            "create_time": "2020-05-15T22:15:51.725+08:00",
            "default": 0
        },
        {
            "bk_module_id": 59,
            "bk_module_name": "m1",
            "create_time": "2020-05-12T21:04:47.286+08:00",
            "default": 0
        },
        {
            "bk_module_id": 58,
            "bk_module_name": "待回收",
            "create_time": "2020-05-12T21:03:37.238+08:00",
            "default": 3
        },
        {
            "bk_module_id": 57,
            "bk_module_name": "故障机",
            "create_time": "2020-05-12T21:03:37.183+08:00",
            "default": 2
        },
        {
            "bk_module_id": 56,
            "bk_module_name": "空闲机",
            "create_time": "2020-05-12T21:03:37.122+08:00",
            "default": 1
        }
    ]
}
```
### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                          |

#### data说明
| 字段      |  类型      |  描述      |
|-----------|------------|------------|
|bk_module_id | int | 模块id |
|bk_module_name | string |模块名称|
|default | int | 表示模块类型 |
|create_time | string | 创建时间 |