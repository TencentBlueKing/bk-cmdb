### 功能描述

查询业务topo树的简要信息，仅包含集群、模块与主机信息。 (v3.9.13)

- 该接口专供GSEKit使用，在ESB文档中为hidden状态

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段            | 类型           | 必选 | 描述                |
|---------------|--------------|----|-------------------|
| bk_biz_id     | int64        | 是  | 业务ID              |
| set_fields    | string array | 是  | 控制返回结果的集群信息里有哪些字段 |
| module_fields | string array | 是  | 控制返回结果的模块信息里有哪些字段 |
| host_fields   | string array | 是  | 控制返回结果的主机信息里有哪些字段 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "set_fields": [
        "bk_set_id",
        "bk_set_name",
        "bk_set_env",
        "bk_platform",
        "bk_system",
        "bk_chn_name",
        "bk_world_id",
        "bk_service_status"
    ],
    "module_fields": [
        "bk_module_id",
        "bk_module_name"
    ],
    "host_fields": [
        "bk_host_id",
        "bk_host_innerip",
        "bk_host_name"
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "set": {
                "bk_set_id": 11,
                "bk_set_name": "set1",
                "bk_set_env": "3",
                "bk_platform": "sq",
                "bk_system": "android",
                "bk_chn_name": "测试集群1",
                "bk_world_id": "35",
                "bk_service_status": "1"
            },
            "modules": [
                {
                    "module": {
                        "bk_module_id": 12,
                        "bk_module_name": "测试模块1"
                    },
                    "hosts": [
                        {
                            "bk_host_id": 13,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "测试主机1"
                        },
                        {
                            "bk_host_id": 23,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "测试主机2"
                        }
                    ]
                },
                {
                    "module": {
                        "bk_module_id": 14,
                        "bk_module_name": "测试模块2"
                    },
                    "hosts": [
                        {
                            "bk_host_id": 15,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "测试主机3"
                        },
                        {
                            "bk_host_id": 24,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "测试主机4"
                        }
                    ]
                }
            ]
        }
    ]
}
```

### 返回结果参数说明

#### response

| 字段         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

#### data

| 字段      | 类型     | 描述   |
|---------|--------|------|
| set     | object | 集群信息 |
| modules | array  | 模块列表 |

#### data.modules

| 字段     | 类型     | 描述   |
|--------|--------|------|
| module | object | 模块信息 |
| hosts  | array  | 主机列表 |

**注意：此处的返回值说明仅对返回值结构进行简单说明，set、module和host具体返回的字段取决于用户自己定义的属性字段**