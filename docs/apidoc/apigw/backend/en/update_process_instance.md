### Description

Batch Update Process Information (Permission: Service Instance Editing Permission)

### Parameters

| Name      | Type  | Required | Description                              |
|-----------|-------|----------|------------------------------------------|
| processes | array | Yes      | Process information, maximum 100 entries |
| bk_biz_id | int   | Yes      | Business ID                              |

#### Explanation of the processes field

| Name                | Type   | Required | Description                    |
|---------------------|--------|----------|--------------------------------|
| process_template_id | int    | No       | Process template ID            |
| auto_start          | bool   | No       | Whether to start automatically |
| auto_time_gap       | int    | No       | Time gap for automatic start   |
| bk_biz_id           | int    | No       | Business ID                    |
| bk_func_id          | string | No       | Function ID                    |
| bk_func_name        | string | No       | Process name                   |
| bk_process_id       | int    | No       | Process ID                     |
| bk_process_name     | string | No       | Process alias                  |
| bk_supplier_account | string | No       | Supplier account               |
| face_stop_cmd       | string | No       | Forced stop command            |
| pid_file            | string | No       | PID file path                  |
| priority            | int    | No       | Startup priority               |
| proc_num            | int    | No       | Number of processes to start   |
| reload_cmd          | string | No       | Process reload command         |
| restart_cmd         | string | No       | Restart command                |
| start_cmd           | string | No       | Start command                  |
| stop_cmd            | string | No       | Stop command                   |
| timeout             | int    | No       | Operation timeout duration     |
| user                | string | No       | Startup user                   |
| work_path           | string | No       | Working directory              |
| process_info        | object | No       | Process information            |

### Request Example

```python
{
  "bk_biz_id": 1,
  "processes": [
    {
      "bk_process_id": 43,
      "description": "",
      "start_cmd": "",
      "restart_cmd": "",
      "pid_file": "",
      "auto_start": false,
      "timeout": 30,
      "auto_time_gap": 60,
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
      "bk_func_id": "",
      "bind_info": [
        {
            "enable": false,  
            "ip": "127.0.0.1",  
            "port": "100",  
            "protocol": "1", 
            "template_row_id": 1  
        }
      ]
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
  "data": [43]
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Data returned by the request                                        |
