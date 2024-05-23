### Description

Query process instance information based on five-segment notation (v3.9.13)

- This interface is specifically designed for GSEKit use and is hidden in the ESB documentation.

### Parameters

| Name             | Type         | Required | Description                                                                                                                                                                                                                          |
|------------------|--------------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id        | int64        | Yes      | Business ID                                                                                                                                                                                                                          |
| bk_set_ids       | int64 array  | No       | Cluster ID list, if empty, it represents any cluster                                                                                                                                                                                 |
| bk_module_ids    | int64 array  | No       | Module ID list, if empty, it represents any module                                                                                                                                                                                   |
| ids              | int64 array  | No       | Service instance ID list, if empty, it represents any instance                                                                                                                                                                       |
| bk_process_names | string array | No       | Process name list, if empty, it represents any process                                                                                                                                                                               |
| bk_process_ids   | int64 array  | No       | Process ID list, if empty, it represents any process                                                                                                                                                                                 |
| fields           | string array | No       | Process attribute list, controls which fields of process instance information are returned, speeding up interface requests and reducing network traffic<br>Empty to return all fields of the process, bk_process_id, bk_process_name |
| page             | dict         | Yes      | Paging conditions                                                                                                                                                                                                                    |

These fields' conditional relationship is AND, only process instances that simultaneously satisfy the filled conditions
will be queried. For example: if both bk_set_ids and bk_module_ids are filled, and bk_module_ids do not belong to
bk_set_ids, the query result will be empty.

#### page

| Name  | Type   | Required | Description                                                                                                           |
|-------|--------|----------|-----------------------------------------------------------------------------------------------------------------------|
| start | int    | No       | Record start position, default is 0                                                                                   |
| limit | int    | Yes      | Number of records per page, maximum is 500                                                                            |
| sort  | string | No       | Sorting field, '-' indicates descending order, can only be a field of the process, default is sorted by bk_process_id |

### Request Example

```json
{
    "bk_biz_id": 3,
    "set": {
        "bk_set_ids": [
            11,
            12
        ]
    },
    "module": {
        "bk_module_ids": [
            60,
            61
        ]
    },
    "service_instance": {
        "ids": [
            4,
            5
        ]
    },
    "process": {
        "bk_process_names": [
            "pr1",
            "alias_pr2"
        ],
        "bk_process_ids": [
            45,
            46,
            47
        ]
    },
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_name"
    ],
    "page": {
        "start": 0,
        "limit": 100,
        "sort": "bk_process_id"
    }
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "count": 2,
        "info": [
            {
                "set": {
                    "bk_set_id": 11,
                    "bk_set_name": "set1",
                    "bk_set_env": "3"
                },
                "module": {
                    "bk_module_id": 60,
                    "bk_module_name": "mm1"
                },
                "host": {
                    "bk_host_id": 4,
                    "bk_cloud_id": 0,
                    "bk_host_innerip": "127.0.0.1",
                    "bk_host_innerip_v6":"1::1",
                    "bk_addressing":"dynamic",
                    "bk_agent_id":"xxxxxx"
                },
                "service_instance": {
                    "id": 4,
                    "name": "127.0.0.1_pr1_3333"
                },
                "process_template": {
                    "id": 48
                },
                "process": {
                    "bk_func_name": "pr1",
                    "bk_process_id": 45,
                    "bk_process_name": "pr1"
                }
            },
            {
                "set": {
                    "bk_set_id": 11,
                    "bk_set_name": "set1",
                    "bk_set_env": "3"
                },
                "module": {
                    "bk_module_id": 60,
                    "bk_module_name": "mm1"
                },
                "host": {
                    "bk_host_id": 4,
                    "bk_cloud_id": 0,
                    "bk_host_innerip": "127.0.0.1"
                },
                "service_instance": {
                    "id": 4,
                    "name": "127.0.0.1_pr1_3333"
                },
                "process_template": {
                    "id": 49
                },
                "process": {
                    "bk_func_name": "pr2",
                    "bk_process_id": 46,
                    "bk_process_name": "alias_pr2"
                }
            }
        ]
    }
}
```

### Response Parameters

| Name    | Type   | Description                                                        |
|---------|--------|--------------------------------------------------------------------|
| result  | bool   | Whether the request is successful. true: successful; false: failed |
| code    | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message | string | Error message returned in case of failure                          |

#### data Field Explanation

| Name             | Type   | Description                                                |
|------------------|--------|------------------------------------------------------------|
| count            | int    | Total number of process instances that meet the conditions |
| set              | object | Cluster information of the process                         |
| module           | object | Module information of the process                          |
| host             | object | Host information of the process                            |
| service_instance | object | Service instance information of the process                |
| process_template | object | Process template information of the process                |
| process          | object | Detailed information of the process itself                 |

#### data.set Field Explanation

| Name        | Type   | Description      |
|-------------|--------|------------------|
| bk_set_id   | int    | Cluster ID       |
| bk_set_name | string | Cluster name     |
| bk_set_env  | string | Environment type |

#### data.module Field Explanation

| Name           | Type   | Description |
|----------------|--------|-------------|
| bk_module_id   | int    | Module ID   |
| bk_module_name | string | Module name |

#### data.host Field Explanation

| Name               | Type   | Description     |
|--------------------|--------|-----------------|
| bk_host_id         | int    | Host ID         |
| bk_cloud_id        | int    | Control area ID |
| bk_host_innerip    | string | Host inner IP   |
| bk_host_innerip_v6 | int    | Host inner IPv6 |
| bk_addressing      | string | Addressing mode |
| bk_agent_id        | string | Agent ID        |

#### data.service_instance Field Explanation

| Name | Type   | Description           |
|------|--------|-----------------------|
| id   | int    | Service instance ID   |
| name | string | Service instance name |

#### data.process_template Field Explanation

| Name | Type | Description         |
|------|------|---------------------|
| id   | int  | Cluster template ID |

#### data.process Field Explanation

| Name                 | Type   | Description                    |
|----------------------|--------|--------------------------------|
| auto_start           | bool   | Whether to automatically start |
| bk_biz_id            | int    | Business ID                    |
| bk_func_name         | string | Process name                   |
| bk_process_id        | int    | Process ID                     |
| bk_process_name      | string | Process alias                  |
| bk_start_param_regex | string | Process startup parameters     |
| bk_supplier_account  | string | Developer account              |
| create_time          | string | Creation time                  |
| description          | string | Description                    |
| face_stop_cmd        | string | Forced stop command            |
| last_time            | string | Update time                    |
| pid_file             | string | PID file path                  |
| priority             | int    | Startup priority               |
| proc_num             | int    | Number of startups             |
| reload_cmd           | string | Process reload command         |
| restart_cmd          | string | Restart command                |
| start_cmd            | string | Start command                  |
| stop_cmd             | string | Stop command                   |
| timeout              | int    | Operation timeout duration     |
| user                 | string | Startup user                   |
| work_path            | string | Working path                   |
| bk_created_at        | string | Creation time                  |
| bk_created_by        | string | Creator                        |
| bk_updated_at        | string | Update time                    |
| bk_updated_by        | string | Updater                        |
| bind_info            | object | Binding information            |
