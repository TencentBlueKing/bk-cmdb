### 描述

根据管控区域名字创建管控区域(权限：管控区域创建权限)

### 输入参数

| 参数名称          | 参数类型   | 必选 | 描述     |
|---------------|--------|----|--------|
| bk_cloud_name | string | 是  | 管控区域名字 |

### 调用示例

```json
{
    
    "bk_cloud_name": "test1"
}

```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "created": {
            "origin_index": 0,
            "id": 6
        }
    }
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

#### data

| 参数名称    | 参数类型   | 描述        |
|---------|--------|-----------|
| created | object | 创建成功，返回信息 |

#### data.created

| 参数名称         | 参数类型 | 描述                  |
|--------------|------|---------------------|
| origin_index | int  | 对应请求的结果顺序           |
| id           | int  | 管控区域id, bk_cloud_id |
