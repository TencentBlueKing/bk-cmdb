### Functional description

Batch create process templates

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id  | int     | yes  | Business ID |
| service_template_id            |  int  |no   | Service template ID|
| processes         |  array  |yes   | Process template information, the max length is 100|


#### processes 
as_default_value: Is the value of the process based on the template

| Field| Type| Required| Description|
|---|---|---|---|
|auto_start| bool| no | Whether to pull up automatically|
|auto_time_gap| int| no | Pull up interval|
|bk_biz_id| int| no | Business ID |
|bk_func_id| string| no | Function ID|
|bk_func_name| string| no | Process name|
|bk_process_id| int| no | Process id|
|bk_process_name| string| no| Process alias|
|bk_supplier_account| string| no| Developer account number|
|face_stop_cmd| string| no| Forced stop command|
|pid_file| string| no| PID file path|
|priority| int| no| Startup priority|
|proc_num| int| no| Number of starts|
|reload_cmd| string| no| Process reload command|
|restart_cmd| string| no| Restart command|
|start_cmd| string| no| Start command|
|stop_cmd| string| no| Stop command|
|timeout| int| no| Operation time-out duration|
|user| string| no| Start user|
|work_path| string| no| Working path|
|bind_info| object| no| Binding information|

#### Bind_info Field Description
| Field| Type| Required| Description|
|---|---|---|---|
|enable| bool| no | Is the port enabled|
|ip| string| no | Bound ip|
|port| string| no | Bound port|
|protocol| string| no | Protocol used|
|row_id| int| no | Template row index used for instantiation, unique in process|

### Request Parameters Example

```json
{ 
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "service_template_id": 1,
  "processes": [
    {
      "spec": {
          "proc_num": {
              "value": null,
              "as_default_value": false
          },
          "stop_cmd": {
              "value": "",
              "as_default_value": false
          },
          "restart_cmd": {
              "value": "",
              "as_default_value": false
          },
          "face_stop_cmd": {
              "value": "",
              "as_default_value": false
          },
          "bk_func_name": {
              "value": "p1",
              "as_default_value": true
          },
          "work_path": {
              "value": "",
              "as_default_value": false
          },
          "priority": {
              "value": null,
              "as_default_value": false
          },
          "reload_cmd": {
              "value": "",
              "as_default_value": false
          },
          "bk_process_name": {
              "value": "p1",
              "as_default_value": true
          },
          "pid_file": {
              "value": "",
              "as_default_value": false
          },
          "auto_start": {
              "value": false,
              "as_default_value": false
          },
          "auto_time_gap": {
              "value": null,
              "as_default_value": false
          },
          "start_cmd": {
              "value": "",
              "as_default_value": false
          },
          "bk_func_id": {
              "value": null,
              "as_default_value": false
          },
          "user": {
              "value": "",
              "as_default_value": false
          },
          "timeout": {
              "value": null,
              "as_default_value": false
          },
          "description": {
              "value": "",
              "as_default_value": false
          },
          "bk_start_param_regex": {
              "value": "",
              "as_default_value": false
          },
          "bind_info": {
              "value": [
                  {
                      "enable": {
                          "value": false,
                          "as_default_value": true
                      },
                      "ip": {
                          "value": "1",
                          "as_default_value": true
                      },
                      "port": {
                          "value": "100",
                          "as_default_value": true
                      },
                      "protocol": {
                          "value": "1",
                          "as_default_value": true
                      },
                      "row_id": 1
                  }
              ],
              "as_default_value": true
          }
        }
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
  "data": [[52]]
}
```

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | array |Successfully created process template ID|
