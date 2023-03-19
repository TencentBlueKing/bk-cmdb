### Functional description

Create service instances in batches. If the module is bound with a service template, the service instances will also be created according to the template. The process template ID corresponding to each process must also be provided in the process parameter for creating the service instance

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_module_id         |  int  |yes   | Module ID|
| instances            | array  | Yes   | new service instance data, the max length is 100 |
| bk_biz_id            |  int  |yes   | Business ID |

#### Instances Field Description

| Field| Type| Required	   | Description| Description|
|---|---|---|---|---|
|instances.bk_host_id| int| yes | Host ID| Host ID of the service instance binding|
|instances.processes| array| yes | Process information| New process information under service instance|
|instances.processes.process_template_id| int| yes | Process template ID| Fill in 0 if the module is not bound to the service template|
|instances.processes.process_info| object| yes | Process instance information| If the process has a template bound to it, only fields in the template that are not locked are valid|

#### Processes Field Description
| Field| Type| Required	   | Description|
|---|---|---|---|
|process_template_id| int| yes | Process template id|
|auto_start| bool| no | Automatically pull up|
|auto_time_gap| int| no | Pull up interval|
|bk_biz_id| int| no | Business ID |
|bk_func_id| string| no | Function ID|
|bk_func_name| string| no | Process name|
|bk_process_id| int| no | Process id|
|bk_process_name| string| no | Process alias|
|bk_supplier_account| string| no | Developer account number|
|face_stop_cmd| string| no | Forced stop command|
|pid_file| string| no | PID file path|
|priority| int| no | Startup priority|
|proc_num| int| no | Number of starts|
|reload_cmd| string| no | Process reload command|
|restart_cmd| string| no | Restart command|
|start_cmd| string| no | Start command|
|stop_cmd| string| no | Stop order|
|timeout| int| no | Operation time-out duration|
|user| string| no | Start user|
|work_path| string| no | Working path|
|process_info| object| yes | Process information|

#### Description of the process_info field
| Field| Type| Required	   | Description|
|---|---|---|---|
|bind_info| object| yes | Binding information|
|bk_supplier_account| string| yes | Developer account number|

#### Bind_info Field Description
| Field| Type| Required	   | Description|
|---|---|---|---|
|enable| bool| yes | Is the port enabled|
|ip| string| yes | Bound ip|
|port| string| yes | Bound port|
|protocol| string| yes | Protocol used|
|template_row_id| int| yes | Template row index used for instantiation, unique in process|

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "bk_module_id": 60,
  "instances": [
    {
      "bk_host_id": 2,
      "processes": [
        {
          "process_template_id": 1,
          "process_info": {
            "bk_supplier_account": "0",
            "bind_info": [
              {
                  "enable": false,
                  "ip": "127.0.0.1",
                  "port": "80",
                  "protocol": "1",
                  "template_row_id": 1234
              }
            ],
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
            "port": "8008,8443",
            "bk_process_name": "job_java",
            "user": "",
            "proc_num": 1,
            "priority": 1,
            "bk_biz_id": 2,
            "bk_func_id": "",
            "bk_process_id": 1
          }
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
  "data": [53]
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
| data | object |New service instance ID list|

