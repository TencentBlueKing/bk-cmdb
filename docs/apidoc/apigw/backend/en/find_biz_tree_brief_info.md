### Description

Query the brief information of the business topology tree, including clusters, modules, and hosts. (v3.9.13)

- This interface is exclusively for GSEKit use and is hidden in the ESB document.

### Parameters

| Name          | Type         | Required | Description                                                  |
|---------------|--------------|----------|--------------------------------------------------------------|
| bk_biz_id     | int64        | Yes      | Business ID                                                  |
| set_fields    | string array | Yes      | Control which fields are included in the cluster information |
| module_fields | string array | Yes      | Control which fields are included in the module information  |
| host_fields   | string array | Yes      | Control which fields are included in the host information    |

### Request Example

```json
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

### Response Example

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
                "bk_chn_name": "Test Cluster 1",
                "bk_world_id": "35",
                "bk_service_status": "1"
            },
            "modules": [
                {
                    "module": {
                        "bk_module_id": 12,
                        "bk_module_name": "Test Module 1"
                    },
                    "hosts": [
                        {
                            "bk_host_id": 13,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "Test Host 1"
                        },
                        {
                            "bk_host_id": 23,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "Test Host 2"
                        }
                    ]
                },
                {
                    "module": {
                        "bk_module_id": 14,
                        "bk_module_name": "Test Module 2"
                    },
                    "hosts": [
                        {
                            "bk_host_id": 15,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "Test Host 3"
                        },
                        {
                            "bk_host_id": 24,
                            "bk_host_innerip": "127.0.0.1",
                            "bk_host_name": "Test Host 4"
                        }
                    ]
                }
            ]
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name    | Type   | Description         |
|---------|--------|---------------------|
| set     | object | Cluster information |
| modules | array  | Module list         |

#### data.modules

| Name   | Type   | Description        |
|--------|--------|--------------------|
| module | object | Module information |
| hosts  | array  | Host list          |

**Note: The return value structure is briefly explained here. The specific fields returned for set, module, and host
depend on the user-defined attribute fields.**
