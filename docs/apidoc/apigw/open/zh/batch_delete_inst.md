### 描述

批量删除对象实例(权限：模型实例删除权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述   |
|-----------|--------|----|------|
| bk_obj_id | string | 是  | 模型ID |
| delete    | object | 是  | 删除   |

#### delete

| 参数名称     | 参数类型  | 必选 | 描述     |
|----------|-------|----|--------|
| inst_ids | array | 是  | 实例ID集合 |

### 调用示例

```json
{
    "bk_obj_id": "bk_firewall",
    "delete": {
        "inst_ids": [
            46,47
        ]
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
    "data": "success"
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
