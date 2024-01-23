### 功能描述

新增模型实例之间的关联关系.(权限：模型实例的编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段              | 类型     | 必选 | 描述            |
|-----------------|--------|----|---------------|
| bk_obj_asst_id  | string | 是  | 模型之间关联关系的唯一id |
| bk_inst_id      | int64  | 是  | 源模型实例id       |
| bk_asst_inst_id | int64  | 是  | 目标模型实例id      |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 11,
    "bk_asst_inst_id": 21
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "id": 1038
    },
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
}

```

### 返回结果参数说明

#### response

| 字段         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |

#### data

| 字段 | 类型    | 描述            |
|----|-------|---------------|
| id | int64 | 新增的实例关联关系身份id |

