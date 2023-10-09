### Functional description

Batch update process information

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| processes            | array  | Yes   | process info, the max length is 100 |
| bk_biz_id            |  int  |yes   | Business ID |

#### Processes Field Description
| Field| Type| Description|
|---|---|---|
|process_template_id| int| Process template id|
|auto_start| bool| Whether to pull up automatically|
|auto_time_gap| int| Pull up interval|
|bk_biz_id| int| Business ID |
|bk_func_id| string| Function ID|
|bk_func_name| string| Process name|
|bk_process_id| int| Process id|
|bk_process_name| string| Process alias|
|bk_supplier_account| string| Developer account number|
|face_stop_cmd| string| Forced stop command|
|pid_file| string| PID file path|
|priority| int| Startup priority|
|proc_num| int| Number of starts|
|reload_cmd| string| Process reload command|
|restart_cmd| string| Restart command|
|start_cmd| string| Start command|
|stop_cmd| string| Stop order|
|timeout| int| Operation time-out duration|
|user| string| Start user|
|work_path| string| Working path|
|process_info| object| Process information|

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "processes": [
    {
      "bk_process_id": 43,
      "bk_supplier_account": "0",
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

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": [43]
}
```

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request succeeded or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

