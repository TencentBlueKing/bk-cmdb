### 描述

新增模型实例之间的关联关系.(权限：模型实例的编辑权限)

### 输入参数

| 参数名称            | 参数类型   | 必选 | 描述            |
|-----------------|--------|----|---------------|
| bk_obj_asst_id  | string | 是  | 模型之间关联关系的唯一id |
| bk_inst_id      | int64  | 是  | 源模型实例id       |
| bk_asst_inst_id | int64  | 是  | 目标模型实例id      |

### 调用示例

```json
{
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 11,
    "bk_asst_inst_id": 21
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "id": 1038
    },
    "permission": null,
}

```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data

| 参数名称 | 参数类型  | 描述            |
|------|-------|---------------|
| id   | int64 | 新增的实例关联关系身份id |
