### Description

Query details of processes corresponding to process IDs under a certain business (v3.9.8)

### Parameters

| Name           | Type  | Required | Description                                                                                                                                                                                                                                                                  |
|----------------|-------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id      | int   | Yes      | Business ID of the process                                                                                                                                                                                                                                                   |
| bk_process_ids | array | Yes      | List of process IDs, up to 500                                                                                                                                                                                                                                               |
| fields         | array | No       | List of process properties, control which fields of process instance information are returned, can speed up the interface request and reduce network traffic transmission <br>When empty, all fields of the process are returned, bk_process_id is a required returned field |

### Request Example

```json
{
    "bk_biz_id":1,
    "bk_process_ids": [
        43,
        44
    ],
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_name"
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
        "data": [
        {
            "auto_start": null,
            "bind_info": [
                {
                    "enable": true,
                    "ip": "127.0.0.1",
                    "port": "9898",
                    "protocol": "1",
                    "template_row_id": 1
                }
            ],
            "bk_biz_id": 3,
            "bk_created_at": "2023-11-15T10:37:39.384+08:00",
            "bk_created_by": "admin",
            "bk_func_name": "delete",
            "bk_process_id": 57,
            "bk_process_name": "delete-aa",
            "bk_start_check_secs": null,
            "bk_start_param_regex": "",
            "bk_supplier_account": "0",
            "bk_updated_at": "2023-11-15T17:19:18.1+08:00",
            "bk_updated_by": "admin",
            "create_time": "2023-11-15T10:37:39.384+08:00",
            "description": "",
            "face_stop_cmd": "",
            "last_time": "2023-11-15T17:19:18.1+08:00",
            "pid_file": "",
            "priority": null,
            "proc_num": null,
            "reload_cmd": "",
            "restart_cmd": "",
            "service_instance_id": 57,
            "start_cmd": "",
            "stop_cmd": "",
            "timeout": null,
            "user": "",
            "work_path": ""
        }
    ],
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | array  | Data returned by the request                                       |

#### data

| Name                 | Type   | Description                    |
|----------------------|--------|--------------------------------|
| auto_start           | bool   | Whether to start automatically |
| bk_biz_id            | int    | Business ID                    |
| bk_func_name         | string | Process name                   |
| bk_process_id        | int    | Process ID                     |
| bk_process_name      | string | Process alias                  |
| bk_start_param_regex | string | Process start parameters       |
| bk_supplier_account  | string | Supplier account               |
| create_time          | string | Creation time                  |
| description          | string | Description                    |
| face_stop_cmd        | string | Forced stop command            |
| last_time            | string | Update time                    |
| pid_file             | string | PID file path                  |
| priority             | int    | Start priority                 |
| proc_num             | int    | Number of starts               |
| reload_cmd           | string | Process reload command         |
| restart_cmd          | string | Restart command                |
| start_cmd            | string | Start command                  |
| stop_cmd             | string | Stop command                   |
| timeout              | int    | Operation timeout duration     |
| user                 | string | Start user                     |
| work_path            | string | Working path                   |
| bk_created_at        | string | Creation time                  |
| bk_created_by        | string | Creator                        |
| bk_updated_at        | string | Update time                    |
| bk_updated_by        | string | Updater                        |
| bind_info            | object | Binding information            |

#### data.info.property.bind_info.value Field Explanation

| Name     | Type   | Description                                   |
|----------|--------|-----------------------------------------------|
| enable   | bool   | Whether the port is enabled                   |
| ip       | string | Bound IP                                      |
| port     | string | Bound port                                    |
| protocol | string | Used protocol                                 |
| row_id   | int    | Template row index, unique within the process |
