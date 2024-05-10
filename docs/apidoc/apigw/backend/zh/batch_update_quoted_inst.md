### 描述

批量更新被引用的模型的实例(v3.10.30+，权限：源模型实例的编辑权限)

### 输入参数

| 参数名称           | 参数类型        | 必选 | 描述                     |
|----------------|-------------|----|------------------------|
| bk_obj_id      | string      | 是  | 源模型ID                  |
| bk_property_id | string      | 是  | 源模型引用该模型的属性ID          |
| ids            | int64 array | 是  | 需要更新的实例ID数组，最多不能超过500个 |
| data           | object      | 是  | 需要更新的实例信息              |

#### data

| 参数名称        | 参数类型   | 必选            | 描述                     |
|-------------|--------|---------------|------------------------|
| name        | string | data中至少一个字段必填 | 名称，此处仅为示例，实际字段由模型属性决定  |
| operator    | string | data中至少一个字段必填 | 维护人，此处仅为示例，实际字段由模型属性决定 |
| description | string | data中至少一个字段必填 | 描述，此处仅为示例，实际字段由模型属性决定  |

### 调用示例

```json
{
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "ids": [
    1,
    2
  ],
  "data": {
    "name": "test",
    "operator": "user",
    "description": "test instance"
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
  "data": null,
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
