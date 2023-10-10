### 功能描述

 查询某实例所有的关联关系（包含其作为关联关系原模型和关联关系目标模型的情况）

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段       | 类型   | 必选 | 描述                  |
| ---------- | ------ | ---- | --------------------- |
| bk_inst_id | int    | 是   | 实例id                |
| bk_obj_id  | string | 是   | 模型id                |
| fields     | array  | 是   | 需要返回的字段        |
| start      | int    | 否   | 记录开始位置          |
| limit      | int    | 否   | 分页大小，最大值500。 |

#### page

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| start     |  int     | 是     | 记录开始位置 |
| limit     |  int     | 是     | 每页限制条数,最大200 |


### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_inst_id": 16,
        "bk_obj_id": "bk_router"
    },
    "fields": [
        "id",
        "bk_inst_id",
        "bk_obj_id",
        "bk_asst_inst_id",
        "bk_asst_obj_id",
        "bk_obj_asst_id",
        "bk_asst_id"
        ],
    "page": {
        "start":0,
        "limit":2
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
    "data": [
        {
            "id": 4,
            "bk_inst_id": 1,
            "bk_obj_id": "bk_switch",
            "bk_asst_inst_id": 16,
            "bk_asst_obj_id": "bk_router",
            "bk_obj_asst_id": "bk_switch_default_bk_router",
            "bk_asst_id": "default"
        },
        {
            "id": 6,
            "bk_inst_id": 2,
            "bk_obj_id": "bk_switch",
            "bk_asst_inst_id": 16,
            "bk_asst_obj_id": "bk_router",
            "bk_obj_asst_id": "bk_switch_default_bk_router",
            "bk_asst_id": "default"
        }
    ]
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

#### data

| 名称            | 类型   | 说明                     |
| :-------------- | :----- | :----------------------- |
| id              | int64  | 关联id                   |
| bk_inst_id      | int64  | 源模型实例id             |
| bk_obj_id       | string | 关联关系源模型id         |
| bk_asst_inst_id | int64  | 关联关系目标模型id       |
| bk_asst_obj_id  | string | 目标模型实例id           |
| bk_obj_asst_id  | string | 自动生成的模型关联关系id |
| bk_asst_id      | string | 关系名称                 |