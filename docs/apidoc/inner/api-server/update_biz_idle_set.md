### 请求方式

POST /api/v3/topo/update/biz/idle_set

### 描述

更改业务空闲机池模块、集群命名信息，若模块为内置模块或已存在模块，则更新命名；否则为新建模块。

### 输入参数

| 参数名称   | 参数类型   | 必选            | 描述                       |
|--------|--------|---------------|--------------------------|
| type   | string | 是             | 更改集群或模块配置(可选值module,set) |
| module | object | module,set二选一 | 更新具体业务空闲机池模块配置信息         |
| set    | object | module,set二选一 | 更新具体业务空闲机池集群配置信息         |

#### module

| 参数名称        | 参数类型   | 必选 | 描述                                       |
|-------------|--------|----|------------------------------------------|
| module_key  | string | 是  | 拓扑key,可选值[内置idle,recycle,fault,用户自定义key] |
| module_name | string | 是  | 业务空闲机池模块名                                |

#### set

| 参数名称     | 参数类型   | 必选 | 描述        |
|----------|--------|----|-----------|
| set_key  | string | 是  | 当前版本支持任意值 |
| set_name | string | 是  | 业务空闲机池集群名 |

### 调用示例
type=module
```json
{
  "type": "module",
  "module": {
    "module_key": "idle",
    "module_name": "空闲模块"
  }
}
```
type=set
```json
{
  "type": "set",
  "set": {
    "set_key": "1",
    "set_name": "空闲机"
  }
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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |
