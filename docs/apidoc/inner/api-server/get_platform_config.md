### 请求方式

GET /api/v3/admin/find/config/platform_config/{type}

### 描述

查询平台配置

### 请求参数

| 参数名称 | 参数类型   | 必选 | 描述                          |
|------|--------|----|-----------------------------|
| type | string | 是  | 可选值：id_generator(id生成器相关配置) |

### 响应示例
type=id_generator
```json
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "permission": null,
  "data": {
    "enabled": false,
    "step": 1,
    "current_id": {
      "biz": 2,
      "host": 0,
      "inst_asst": 0,
      "module": 5,
      "object_instance": 0,
      "process": 0,
      "service_instance": 0,
      "set": 2
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

#### data说明
type=id_generator

| 参数名称       | 参数类型           | 描述                 |
|------------|----------------|--------------------|
| enabled    | bool           | 是否开启id generator配置 |
| step       | int            | id generator配置步长   |
| current_id | map[string]int | 模型实例当前id           |
