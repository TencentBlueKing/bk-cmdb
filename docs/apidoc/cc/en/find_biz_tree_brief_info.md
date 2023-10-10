### Functional description
Query the brief information of the service topo tree, which only contains  set , module and host information. (v3.9.13)

- This interface is intended for use by GSEKit and is hidden in the ESB documentation

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

|Field| Type| Required| Description|
|---|---|---|---|
|bk_biz_id| int64| yes | Business ID |
|set_fields| string array| yes | Controls which fields are in the set information that returns the result |
|module_fields| string array| yes | Controls which fields are in the module information that returns the result|
|host_fields| string array| yes | Controls which fields are in the host information that returns the result|


### Request Parameters Example

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

### Return Result Example
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