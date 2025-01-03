### Function Description

Query process instance list based on service instance ID

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type | Required | Description         |
| ------------------- | ---- | -------- | ------------------- |
| bk_biz_id           | int  | Yes      | Business ID         |
| service_instance_id | int  | Yes      | Service instance ID |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "service_instance_id": 54
}
```

### Response Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": [
    {
      "property": {
        "auto_start": false,
        "auto_time_gap": 60,
        "bk_biz_id": 2,
        "bk_func_id": "",
        "bk_func_name": "java",
        "bk_process_id": 46,
        "bk_process_name": "job_java",
        "bk_start_param_regex": "",
        "create_time": "2019-06-05T14:59:12.065+08:00",
        "description": "",
        "face_stop_cmd": "",
        "last_time": "2019-06-05T14:59:12.065+08:00",
        "pid_file": "",
        "priority": 1,
        "proc_num": 1,
        "reload_cmd": "",
        "restart_cmd": "",
        "start_cmd": "",
        "stop_cmd": "",
        "timeout": 30,
        "user": "",
        "work_path": "/data/bkee",
        "bind_info": [
            {
                "enable": false,  
                "ip": "127.0.0.1",  
                "port": "100",  
                "protocol": "1", 
                "template_row_id": 1  
            }
        ]
      },
      "relation": {
        "bk_biz_id": 1,
        "bk_process_id": 46,
        "service_instance_id": 54,
        "process_template_id": 1,
        "bk_host_id": 1,
      }
    }
  ]
}
```

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | array  | Data returned by the request                                 |

#### data Field Explanation

| Field    | Type   | Description                                       |
| -------- | ------ | ------------------------------------------------- |
| property | object | Process property information                      |
| relation | object | Relationship between process and service instance |

#### data[x].property Field Explanation

| Field                | Type   | Description                    |
| -------------------- | ------ | ------------------------------ |
| auto_start           | bool   | Whether to start automatically |
| bk_biz_id            | int    | Business ID                    |
| bk_func_name         | string | Process name                   |
| bk_process_id        | int    | Process ID                     |
| bk_process_name      | string | Process alias                  |
| bk_start_param_regex | string | Process start parameters       |
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

#### data[x].property.bind_info[n] Field Explanation

| Field           | Type   | Description                                   |
| --------------- | ------ | --------------------------------------------- |
| enable          | bool   | Whether the port is enabled                   |
| ip              | string | Bound IP                                      |
| port            | string | Bound port                                    |
| protocol        | string | Used protocol                                 |
| template_row_id | int    | Template row index, unique within the process |

#### data[x].relation Field Explanation

| Field               | Type   | Description         |
| ------------------- | ------ | ------------------- |
| bk_biz_id           | int    | Business ID         |
| bk_process_id       | int    | Process ID          |
| service_instance_id | int    | Service instance ID |
| process_template_id | int    | Process template ID |
| bk_host_id          | int    | Host ID             |
