### 描述

根据业务id和服务实例id，以及需要移除的标签，将指定业务下的服务实例移除标签(权限：服务实例删除权限)

### 输入参数

| 参数名称         | 参数类型  | 必选 | 描述               |
|--------------|-------|----|------------------|
| bk_biz_id    | int   | 是  | 业务ID             |
| instance_ids | array | 是  | 服务实例ID列表,最大值为500 |
| keys         | array | 是  | 需要移除的标签key列表     |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "instance_ids": [60, 62],
  "keys": ["value1", "value3"]
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
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
