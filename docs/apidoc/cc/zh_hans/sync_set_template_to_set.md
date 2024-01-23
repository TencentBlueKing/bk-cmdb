### 功能描述

根据业务id，集群模板id,待同步集群id列表，将指定业务下的集群模板同步到集群(权限：集群编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段              | 类型    | 必选 | 描述        |
|-----------------|-------|----|-----------|
| bk_biz_id       | int   | 是  | 业务ID      |
| set_template_id | int   | 是  | 集群模板ID    |
| bk_set_ids      | array | 是  | 待同步集群ID列表 |

### 请求参数示例

```json
{

    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 20,
    "set_template_id": 6,
    "bk_set_ids": [46]
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

| 字段         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |
