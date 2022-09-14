### Functional description

Update process template information

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| process_template_id            |  int  |no   | Process template ID|
| process_property         |  object  |yes   | Process template field information to update|

#### The fields where process_property can appear

Note:

as_default_value: Is the value of the process based on the template
Value: the value of the process. Different field types are different

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
|proc_num|  object| no |{"value": null, "as_default_value": False}, value type is numeric|
|stop_cmd| object| no |{"value": "","as_default_value": False}, the value type is a string|
|restart_cmd| object| no |{"value": "","as_default_value": False}, the value type is a string|
|face_stop_cmd| object| no |{"value": "","as_default_value": False}, the value type is a string|
|bk_func_name| object| no |{"value": "a7","as_default_value": True}}, value type is string|
|work_path| object| no |{"value": "","as_default_value": False}, the value type is a string|
|priority| object| no |{"value": null,"as_default_value": False}, value type is numeric|
|reload_cmd| object| no |{"value": "","as_default_value": False}, the value type is a string|
|bk_process_name| object| no |{"value": "a7","as_default_value": True}}, value type is string|
|pid_file| object| no |{"value": "","as_default_value": False}, value type is a string|
|auto_start| object| no |{"value": null,"as_default_value": Null}}, value type is boolean|
|auto_time_gap| object| no |{"value": null,"as_default_value": False}, value type is numeric|
|start_cmd| object| no |{"value": "","as_default_value": False}, the value type is a string|
|bk_func_id| object| no |{"value": "","as_default_value": False} the value type is a string|
|user| object| no |{"value": "","as_default_value": False}, the value type is a string|
|timeout| object| no |{"value": null,"as_default_value": False}, value type is numeric|
|description| object| no |{"value": "1","as_default_value": True}}, value type is string|
|bk_start_param_regex| object| no |{"value": "","as_default_value": False}, the value type is a string|
|bind_info| object| no |{"value":[],,"as_default_value": True }, see access_property.bind_info.value n for details of value[]|


#### Process_property.bind_info.value [n] fields that can appear

Note:

When modifying bind_info, you must first obtain the bind_info content of the original process, then modify it on the existing bind_info of the process, and transfer the modified content to the modification structure.

Note:

as_default_value: Is the value of the process based on the template
Value: value of process. Different field types are different

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
|enable| object| no |{"value": false,"as_default_value": True}, value type is boolean|
|ip| object| no |{"value": "1","as_default_value": True}, the value type is a string|
|port| object| no |{"value": "100","as_default_value": True}, the value type is a string|
|protocol| object| no |{"value": "1","as_default_value": True},, value type is a string|
|row_id| int| no | Unique representation id, new row can be set to empty, update must keep the original value|







### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "process_template_id": 50,
  "process_property": {
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
  }
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
    "id": 50,
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
| result | bool |Whether the request succeeded or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Updated process template information|

#### Data field Description

| Name| Type| Description|
|---|---|---|
| id | int |Data id|
| bk_process_name | string |Process alias|
| bk_biz_id |  int| Business ID |
| service_template_id | int |Service template id|
| property |object  |Attribute|
| creator              |  string             | Creator of this data                                                                                 |
| modifier             |  string             | The last person to modify this piece of data            |
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string       | Developer account number|
