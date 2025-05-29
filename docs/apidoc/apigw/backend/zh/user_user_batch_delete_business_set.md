### 描述

删除业务集(版本：v3.10.12+，权限：业务集删除权限)

### 输入参数

| 参数名称           | 参数类型  | 必选 | 描述      |
|----------------|-------|----|---------|
| bk_biz_set_ids | array | 是  | 业务集ID列表 |

### 调用示例

```json
{
    "bk_biz_set_ids":[
        10,
        12
    ]
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission":null,
    "data": {},
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
