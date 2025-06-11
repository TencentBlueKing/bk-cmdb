### 描述

对指定业务id下通过集群id批量删除集群(权限：业务拓扑删除权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述   |
|-----------|--------|----|------|
| bk_biz_id | int    | 是  | 业务ID |
| delete    | object | 是  | 删除   |

#### delete

| 参数名称     | 参数类型      | 必选 | 描述     |
|----------|-----------|----|--------|
| inst_ids | int array | 是  | 集群ID集合 |

### 调用示例

```json
{
    "bk_biz_id":0,
    "delete": {
    "inst_ids": [123]
    }
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
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
