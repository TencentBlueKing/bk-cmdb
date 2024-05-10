### 描述

根据业务ID获取业务空闲机, 故障机和待回收模块

### 输入参数

| 参数名称      | 参数类型 | 必选 | 描述   |
|-----------|------|----|------|
| bk_biz_id | int  | 是  | 业务ID |

### 调用示例

```json
{
    "bk_biz_id":0
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
    "bk_set_id": 2,
    "bk_set_name": "idle pool",
    "module": [
      {
        "bk_module_id": 3,
        "bk_module_name": "idle host",
        "default": 1,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 4,
        "bk_module_name": "fault host",
        "default": 2,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 5,
        "bk_module_name": "recycle host",
        "default": 3,
        "host_apply_enabled": false
      }
    ]
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

| 参数名称        | 参数类型   | 描述                        |
|-------------|--------|---------------------------|
| bk_set_id   | int64  | 空闲机, 故障机和待回收模块所属的set的实例ID |
| bk_set_name | string | 空闲机, 故障机和待回收模块所属的set的实例名称 |
| module      | array  | 空闲机, 故障机和待回收模块信息          |

#### module说明

| 参数名称               | 参数类型   | 描述                  |
|--------------------|--------|---------------------|
| bk_module_id       | int    | 空闲机, 故障机或待回收模块的实例ID |
| bk_module_name     | string | 空闲机, 故障机或待回收模块的实例名称 |
| default            | int    | 表示模块类型              |
| host_apply_enabled | bool   | 是否启用主机属性自动应用        |
