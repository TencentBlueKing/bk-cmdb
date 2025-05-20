### 请求方式

GET /api/v3/admin/find/config/global_config

### 描述

查询全局配置，包括拓扑层级、校验规则、主机机池及业务空闲机池集群、模块配置信息

### 响应示例

```json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "permission": null,
    "data": {
        "backend": {
            "max_biz_topo_level": 7
        },
        "create_time": "2025-05-09T13:05:20.482Z",
        "idle_pool": {
            "fault": "故障机",
            "idle": "空闲机",
            "recycle": "待回收",
            "user_modules": null
        },
        "last_time": "2025-05-09T13:05:40.482Z",
        "set": "空闲机池",
        "validation_rules": {
            "associationId": {
                "description": "关联类型唯一标识验证规则",
                "i18n": {
                    "cn": "由英文字符开头，和下划线、数字或英文组合的字符",
                    "en": "Start with lowercase or uppercase letter, followed by lowercase / uppercase / underscore / numbers characters"
                },
                "value": "XlthLXpBLVpdW1x3XSok"
            }
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

| 参数名称             | 参数类型              | 描述         |
|------------------|-------------------|------------|
| backend          | object            | 拓扑层级配置     |
| idle_pool        | object            | 业务空闲机池模块配置 |
| set              | string            | 业务空闲机池集群配置 |
| validation_rules | map[string]object | 校验规则       |

#### data.backend

| 参数名称               | 参数类型  | 描述       |
|--------------------|-------|----------|
| max_biz_topo_level | int64 | 拓扑最大可建层级 |

#### data.idle_pool

| 参数名称         | 参数类型              | 描述      |
|--------------|-------------------|---------|
| idle         | string            | 空闲模块命名  |
| fault        | string            | 故障模块命名  |
| recycle      | string            | 待回收模块命名 |
| user_modules | map[string]string | 用户自定义模块 |

#### data.validation_rules

| 参数名称         | 参数类型              | 描述        |
|--------------|-------------------|-----------|
| idle         | string            | 空闲模块命名    |
| fault        | string            | 地址，可以填写多个 |
| recycle      | string            | 待回收模块命名   |
| user_modules | map[string]string | 用户自定义模块   |

