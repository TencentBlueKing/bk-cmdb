### 功能描述

根据进程模板ID列表删除进程模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id           | int        | 是     | 业务id     |
| process_templates | array  | 是   | 进程模板ID列表,最大值为500 |

### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "process_templates": [50]
}
```

### 返回结果示例

```python
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "069cbd4eed2846a0b4c995f3d040e2a5"
}
```

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |
