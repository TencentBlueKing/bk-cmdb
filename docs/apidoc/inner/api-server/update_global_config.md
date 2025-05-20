### 请求方式

PUT /api/v3/admin/update/config/global_config/{type}

### 描述

更新全局配置

### 请求参数

| 参数名称 | 参数类型   | 必选 | 描述                    |
|------|--------|----|-----------------------|
| type | string | 是  | 可选值:backend(更新拓扑层级配置) |

### 输入参数
type=backend

| 参数名称            | 参数类型              | 必选 | 描述                 |
|-----------------|-------------------|----|--------------------|
| backend         | object            | 是  | 修改全局配置的拓扑层级配置      |

#### backend

| 参数名称               | 参数类型 | 必选 | 描述   |
|--------------------|------|----|------|
| max_biz_topo_level | int  | 是  | 拓扑层级 |

### 调用示例
type=backend
```json
{
  "backend": {
    "max_biz_topo_level": 10
  }
}
```

### 响应示例
type=backend
```json
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "permission": null,
  "data": "update general backend config success"
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
