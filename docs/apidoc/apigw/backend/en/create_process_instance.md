### Description

Create Process Instance Based on Service Instance ID and Process Instance Attribute Values (Permission: Service Instance
Edit Permission)

### Parameters

| Name                | Type  | Required | Description                                                                             |
|---------------------|-------|----------|-----------------------------------------------------------------------------------------|
| service_instance_id | int   | Yes      | Service instance ID                                                                     |
| processes           | array | Yes      | Attribute values in process instance different from the template, with a maximum of 100 |

#### Explanation of process_info Fields

| Name                | Type   | Required | Description                       |
|---------------------|--------|----------|-----------------------------------|
| auto_start          | bool   | No       | Whether to start automatically    |
| auto_time_gap       | int    | No       | Time interval for automatic start |
| bk_biz_id           | int    | No       | Business ID                       |
| bk_func_id          | string | No       | Function ID                       |
| bk_func_name        | string | No       | Process name                      |
| bk_process_id       | int    | No       | Process ID                        |
| bk_process_name     | string | No       | Process alias                     |
| bk_supplier_account | string | No       | Developer account                 |
| face_stop_cmd       | string | No       | Forced stop command               |
| pid_file            | string | No       | PID file path                     |
| priority            | int    | No       | Startup priority                  |
| proc_num            | int    | No       | Number of startups                |
| reload_cmd          | string | No       | Process reload command            |
| restart_cmd         | string | No       | Restart command                   |
| start_cmd           | string | No       | Start command                     |
| stop_cmd            | string | No       | Stop command                      |
| timeout             | int    | No       | Operation timeout duration        |
| user                | string | No       | Startup user                      |
| work_path           | string | No       | Working directory                 |
| bind_info           | object | No       | Binding information               |

#### Explanation of bind_info Fields

| Name     | Type   | Required | Description                                                          |
|----------|--------|----------|----------------------------------------------------------------------|
| enable   | bool   | No       | Whether the port is enabled                                          |
| ip       | string | No       | Bound IP                                                             |
| port     | string | No       | Bound port                                                           |
| protocol | string | No       | Used protocol                                                        |
| row_id   | int    | No       | Template row index used for instantiation, unique within the process |

### Request Example

```json
{
  "bk_biz_id": 1,
  "service_instance_id": 48,
  "processes": [
    {
      "process_info": {
        "description": "",
        "start_cmd": "",
        "restart_cmd": "",
        "pid_file": "",
        "auto_start": false,
        "timeout": 30,
        "reload_cmd": "",
        "bk_func_name": "java",
        "work_path": "/data/bkee",
        "stop_cmd": "",
        "face_stop_cmd": "",
        "bk_process_name": "job_java",
        "user": "",
        "proc_num": 1,
        "priority": 1,
        "bk_biz_id": 2,
        "bk_start_param_regex": "",
        "bk_process_id": 1,
        "bind_info": [
          {
              "enable": false,
              "ip": "127.0.0.1",
              "port": "80",
              "protocol": "1",
              "template_row_id": 1234
          }
        ]
      }
    }
  ]
}
```

### Response Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [64]
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Newly created process instance ID list                                      |
