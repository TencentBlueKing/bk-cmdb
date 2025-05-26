### 请求方式

PUT /api/v3/admin/update/config/platform_config/{type}

### 描述

更新平台配置

### 请求参数

| 参数名称 | 参数类型   | 必选 | 描述                                 |
|------|--------|----|------------------------------------|
| type | string | 是  | 可选值id_generator 查询id_generator配置内容 |

### 输入参数
type=id_generator

| 参数名称         | 参数类型   | 必选 | 描述                 |
|--------------|--------|----|--------------------|
| id_generator | object | 是  | 更新id_generator配置内容 |

#### id_generator

| 参数名称    | 参数类型   | 必选 | 描述                 |
|---------|--------|----|--------------------|
| step    | int    | 是  | id generator步长     |
| enabled | string | 是  | 是否开始id generator配置 |
| init_id | int    | 否  | 更新模型实例id的初始配置值     |

### 调用示例

[type=id_generator]
```json
{
    "id_generator": {
        "enabled": true,
        "step": 2
    },
  "init_id": {
    "biz": 4
  }
}
```

### 响应示例

[type=id_generator]
```json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "permission": null,
    "data": "update platform config id_generator success"
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
