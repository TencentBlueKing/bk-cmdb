### 功能描述

批量删除被引用的模型的实例(v3.10.30+，权限：源模型实例的编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段             | 类型          | 必选  | 描述                     |
|----------------|-------------|-----|------------------------|
| bk_obj_id      | string      | 是   | 源模型ID                  |
| bk_property_id | string      | 是   | 源模型引用该模型的属性ID          |
| ids            | int64 array | 是   | 需要删除的实例ID数组，最多不能超过500个 |

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "ids": [
    1,
    2
  ]
}
```

### 返回结果示例

```json

{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null,
  "request_id": "dsda1122adasadadada2222"
}
```

### 返回结果参数说明

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |
