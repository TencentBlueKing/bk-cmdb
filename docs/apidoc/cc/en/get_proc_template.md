### Functional description

Obtain process template information, and specify process template ID in url parameter

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            |  int        | yes  | Business ID |
| process_template_id  | int        | yes  | Process template ID     |

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "process_template_id": 49,
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "id": 49,
    "bk_process_name": "p1",
    "bk_biz_id": 1,
    "service_template_id": 51,
    "property": {
      "proc_num": {
        "value": 300,
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
        "value": "root100",
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
    },
    "creator": "admin",
    "modifier": "admin",
    "create_time": "2019-06-19T15:24:04.763+08:00",
    "last_time": "2019-06-21T16:25:03.962512+08:00",
    "bk_supplier_account": "0"
  }
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
| data | object |Data returned by request|

#### Data field Description
| Name| Type| Description|
|---|---|---|
| id |int  |Process template id|
| bk_process_name |string  |Process alias|
| bk_biz_id |  int| Business ID |
|  service_template_id|   int| Service template id|
| property | object |Process properties|
| creator              |  string             | Creator of this data                                                                                 |
| modifier             |  string             | The last person to modify this piece of data            |
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string       | Developer account number|

#### Property Field Description

| Field| Type| Description|
|---|---|---|
|auto_start| bool| Whether to pull up automatically|
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

#### Bind_info Field Description
| Field| Type| Description|
|---|---|---|
|enable| bool| Is the port enabled|
|ip| string| Bound ip|
|port| string| Bound port|
|protocol| string| Protocol used|
|row_id| int| Template row index used for instantiation, unique in process|