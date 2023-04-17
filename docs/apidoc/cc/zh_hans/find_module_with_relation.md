### 功能描述

根据条件查询业务下的模块 (v3.9.7)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  | 是     | 业务ID |
| bk_set_ids  |  array  | 否     | 集群ID列表, 最多可填200个 |
| bk_service_template_ids  |  array  | 否     | 服务模板ID列表 |
| fields  |   array   | 是     | 模块属性列表，控制返回结果的模块信息里有哪些字段 |
| page       |  object    | 是     | 分页信息 |

#### page 字段说明

| 字段  | 类型   | 必选 | 描述                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | 是   | 记录开始位置          |
| limit | int    | 是   | 每页限制条数,最大500 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 2,
    "bk_set_ids":[1,2],
    "bk_service_template_ids": [3,4],
    "fields":["bk_module_id", "bk_module_name"],
    "page": {
        "start": 0,
        "limit": 10
    }
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
    "data": {
        "count": 2,
        "info": [
            {
                "bk_module_id": 8,
                "bk_module_name": "license"
            },
            {
                "bk_module_id": 12,
                "bk_module_name": "gse_proc"
            }
        ]
    }
}
```
### 返回结果参数说明
#### response
| 名称    | 类型   | 说明                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

data 字段说明：

| 名称     | 类型         | 说明               |
| -------- | ------------ | ------------------ |
| count    | int          | 记录条数           |
| info | object array | 模块实际数据 |