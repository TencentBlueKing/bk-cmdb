### 功能描述

根据模型实例关联关系的唯一身份id,删除模型实例之间的关联关系。(权限：源模型实例和目标模型实例的编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数
| 字段                 |  类型      | 必填	   |  描述          |
|----------------------|------------|--------|-----------------------------|
| id           | int     | 是    | 模型实例关联关系的唯一身份id             |
| bk_obj_id    | string    | 是    | 模型实例关联关系的源或目标模型id(v3.10+) |

### 请求参数示例

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "test",
    "id": 1
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
    "data": null
}

```

### 返回结果参数说明

#### response

| 字段       | 类型     | 描述         |
|------------|----------|--------------|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | object | 请求返回的数据 |

