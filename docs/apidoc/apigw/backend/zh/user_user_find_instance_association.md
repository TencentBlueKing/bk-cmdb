### 描述

查询模型的实例关联关系。(权限：模型实例查询权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述            |
|-----------|--------|----|---------------|
| condition | object | 是  | 查询条件          |
| bk_obj_id | string | 是  | 源模型id(v3.10+) |

#### condition

| 参数名称           | 参数类型   | 必选 | 描述          |
|----------------|--------|----|-------------|
| bk_obj_asst_id | string | 是  | 模型关联关系的唯一id |
| bk_asst_id     | string | 否  | 关联类型的唯一id   |
| bk_asst_obj_id | string | 否  | 目标模型id      |

### 调用示例

```json
{
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_asst_obj_id": ""
    },
    "bk_obj_id": "xxx"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [{
        "id": 481,
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
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

| 参数名称            | 参数类型   | 描述                          |
|-----------------|--------|-----------------------------|
| id              | int    | the association's unique id |
| bk_obj_asst_id  | string | 自动生成的模型关联关系id.              |
| bk_obj_id       | string | 关联关系源模型id                   |
| bk_asst_obj_id  | string | 关联关系目标模型id                  |
| bk_inst_id      | int    | 源模型实例id                     |
| bk_asst_inst_id | int    | 目标模型实例id                    |
