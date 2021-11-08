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

### 请求参数示例

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

### 返回结果示例

```json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
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

### 返回结果参数说明

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