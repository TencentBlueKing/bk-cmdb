### 描述

查询某实例所有的关联关系（包含其作为关联关系原模型和关联关系目标模型的情况，权限：模型实例查询权限）

### 输入参数

| 参数名称       | 参数类型   | 必选 | 描述           |
|------------|--------|----|--------------|
| bk_inst_id | int    | 是  | 实例id         |
| bk_obj_id  | string | 是  | 模型id         |
| fields     | array  | 是  | 需要返回的字段      |
| start      | int    | 否  | 记录开始位置       |
| limit      | int    | 否  | 分页大小，最大值500。 |

#### page

| 参数名称  | 参数类型 | 必选 | 描述           |
|-------|------|----|--------------|
| start | int  | 是  | 记录开始位置       |
| limit | int  | 是  | 每页限制条数,最大200 |

### 调用示例

```json
{
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

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称            | 参数类型   | 描述            |
|-----------------|--------|---------------|
| id              | int64  | 关联id          |
| bk_inst_id      | int64  | 源模型实例id       |
| bk_obj_id       | string | 关联关系源模型id     |
| bk_asst_inst_id | int64  | 关联关系目标模型id    |
| bk_asst_obj_id  | string | 目标模型实例id      |
| bk_obj_asst_id  | string | 自动生成的模型关联关系id |
| bk_asst_id      | string | 关系名称          |
