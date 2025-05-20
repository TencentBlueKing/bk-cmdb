### 请求方式

POST /api/v3/topo/delete/biz/extra_module

### 描述

全局配置，删除所有业务中空闲机池中用户自定义拓扑模块

### 请求参数

| 参数名称        | 参数类型   | 必选 | 描述         |
|-------------|--------|----|------------|
| module_key  | string | 是  | 用户自定义模块key |
| module_name | string | 是  | 用户自定义模块名   |

### 调用示例

```json
{
  "module_key": "module_key",
  "module_name": "module_name"
}
```


### 响应示例

```json
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "permission": null,
  "data": null
}
```

