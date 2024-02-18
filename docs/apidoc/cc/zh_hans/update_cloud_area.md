### 功能描述

更新管控区域(权限：管控区域编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段            | 类型     | 必选 | 描述     |
|---------------|--------|----|--------|
| bk_cloud_id   | int    | 是  | 管控区域ID |
| bk_cloud_name | string | 否  | 管控区域名称 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_cloud_id": 5,
    "bk_cloud_name": "管控区域1"
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
    "data": "success"
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
