### 功能描述
查询业务topo树的简要信息，仅包含集群、模块与主机信息。 (v3.9.13)

- 该接口专供GSEKit使用，在ESB文档中为hidden状态

### 请求参数

{{ common_args_desc }}

#### 接口参数

|字段|类型|必填|描述|
|---|---|---|---|
|bk_biz_id|int64|Yes|业务ID|
|set_fields|string array|Yes|控制返回结果的集群信息里有哪些字段|
|module_fields|string array|Yes|控制返回结果的模块信息里有哪些字段|
|host_fields|string array|Yes|控制返回结果的主机信息里有哪些字段|


### 请求参数示例

``` json
{
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
``` json
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