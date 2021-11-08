### Functional description

update process template info

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| process_template_id            | int  | Yes   | process template id |
| process_property         | object  | Yes   | process template property |


#### Process_property fields that can appear

annotation:

as_default_value: Whether the value of the process is based on the template
value: the value of the process, different field types are different

| Field | Type | Required | Description |
|----------------------|------------|--------|-----------------------|
|proc_num| object| no| {"value": null, "as_default_value": false}, value type is number|
|stop_cmd|object| no| {"value": "","as_default_value": false}, value type is a string|
|restart_cmd|object|no|{"value": "","as_default_value": false}, the value type is a string|
|face_stop_cmd|object|no|{"value": "","as_default_value": false}, the value type is a string|
|bk_func_name|object|no|{"value": "a7","as_default_value": true}}, the value type is a string|
|work_path|object|no|{"value": "","as_default_value": false}, the value type is a string|
|priority|object|no|{"value": null,"as_default_value": false}, value type is number|
|reload_cmd|object|no|{"value": "","as_default_value": false}, the value type is a string|
|bk_process_name|object|no|{"value": "a7","as_default_value": true}}, the value type is a string|
|pid_file|object|no|{"value": "","as_default_value": false}, the value type is a string|
|auto_start|object|no|{"value": null,"as_default_value": null}}, value type is boolean|
|auto_time_gap|object|no|{"value": null,"as_default_value": false}, value type is number|
|start_cmd|object|no|{"value": "","as_default_value": false}, the value type is a string|
|bk_func_id|object|no|{"value": "","as_default_value": false} The value type is a string|
|user|object|no|{"value": "","as_default_value": false}, the value type is a string|
|timeout|object|no|{"value": null,"as_default_value": false}, value type is number|
|description|object|no|{"value": "1","as_default_value": true}}, the value type is a string|
|bk_start_param_regex|object|no|{"value": "","as_default_value": false}, the value type is a string||
|bind_info|object|no| {"value":[],,"as_default_value": true }, for value details see process_property.bind_info.value[n]|


#### process_property.bind_info.value[n] Fields that can appear


note:

When modifying bind_info, you must first obtain the content of bind_info of the original process, then modify the existing bind_info of the process, and pass the modified content to the modified structure.

annotation:

as_default_value: Whether the value of the process is based on the template
value: the value of the process, different field types are different

| Field | Type | Required | Description |
|----------------------|------------|--------|-----------------------|
|enable|object|no| {"value": false,"as_default_value": true}, value type is boolean|
|ip|object|no| {"value": "1","as_default_value": true}, value type is string||
|port|object|no| {"value": "100","as_default_value": true}, value type is string||
|protocol|object|no| {"value": "1","as_default_value": true}, value type is string||
|row_id|int|no| the unique id, the newly added row can be set to empty, and the update must keep the original value|


### Request Parameters Example

```python
{
  "bk_biz_id": 1,
  "process_template_id": 50,
  "process_property": {
    "user": {
      "as_default_value": true,
      "value": "root100"
    },
    "proc_num": {
      "as_default_value": true,
      "value": 300
    }
  }
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
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

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |
