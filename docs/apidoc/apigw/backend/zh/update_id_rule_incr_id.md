### 描述

更新id规则自增id (版本：v3.14.1，权限：id规则自增id编辑权限)

### 输入参数

| 参数名称        | 参数类型   | 必选 | 描述                                                  |
|-------------|--------|----|-----------------------------------------------------|
| type        | string | 是  | 如果是更新模型的自增id，则为对应模型的bk_obj_id;如果是全局的自增id，则为"global" |
| sequence_id | int    | 是  | 自增ID                                                |

### 调用示例

```json
{
    "type": "host",
    "sequence_id": 1000000
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": null
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
