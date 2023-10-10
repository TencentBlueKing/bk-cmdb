### Functional description

Query process instance list based on service instance ID

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
|bk_biz_id| int| yes | Business ID |
| service_instance_id | int  |yes   | Service instance ID|


### Request Parameters Example

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

### Return Result Example

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
        "bk_supplier_account": "0",
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
        "bk_supplier_account": ""
      }
    }
  ]
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
| data | array |Data returned by request|

#### Data field Description

| Field| Type| Description|
|---|---|---|
|property| object| Process attribute information|
|relation| object| Process and service instance association information|

#### data [x]. Property Field Description
| Field| Type| Description|
|---|---|---|
|auto_start| bool| Automatically pull up|
|auto_time_gap| int| Pull up interval|
|bk_biz_id| int| Business ID |
|bk_func_id| string| Function ID|
|bk_func_name| string| Process name|
|bk_process_id| int| Process id|
|bk_process_name| string| Process alias|
|bk_start_param_regex| string| Process start parameters|
|bk_supplier_account| string| Developer account number|
|create_time| string| Settling time|
|description| string| Description|
|face_stop_cmd| string| Forced stop command|
|last_time| string| Update time|
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
|bind_info| object| Binding information|

#### data [x] .Property.bind .Property.bind info [n] Field Description
| Field| Type| Description|
|---|---|---|
|enable| bool| Is the port enabled|
|ip| string| Bound ip|
|port| string| Bound port|
|protocol| string| Protocol used|
|template_row_id| int| Template row index used for instantiation, unique in process|

#### data [x]. Recall Field Description
| Field| Type| Description|
|---|---|---|
|bk_biz_id| int| Business ID |
|bk_process_id| int| Process id|
|service_instance_id| int| Service instance id|
|process_template_id| int| Process template id|
|bk_host_id| int| Host id|
|bk_supplier_account| string| Developer account number|