### Description

Update Process Template Information (Permission: Service Template Editing Permission)

### Parameters

| Name                | Type   | Required | Description                                                           |
|---------------------|--------|----------|-----------------------------------------------------------------------|
| process_template_id | int    | Yes      | Process template ID                                                   |
| process_property    | object | Yes      | Information of fields in the process template that need to be updated |

#### Fields that can appear in process_property

Note:

as_default_value: Whether the value of the process is based on the template value: The value of the process, different
field types have different types

| Name                 | Type   | Required | Description                                                                                                |
|----------------------|--------|----------|------------------------------------------------------------------------------------------------------------|
| proc_num             | object | No       | {"value": null, "as_default_value": false}, value type is number                                           |
| stop_cmd             | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| restart_cmd          | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| face_stop_cmd        | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| bk_func_name         | object | No       | {"value": "a7","as_default_value": true}}, value type is string                                            |
| work_path            | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| priority             | object | No       | {"value": null,"as_default_value": false}, value type is number                                            |
| reload_cmd           | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| bk_process_name      | object | No       | {"value": "a7","as_default_value": true}}, value type is string                                            |
| pid_file             | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| auto_start           | object | No       | {"value": null,"as_default_value": null}}, value type is boolean                                           |
| auto_time_gap        | object | No       | {"value": null,"as_default_value": false}, value type is number                                            |
| start_cmd            | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| bk_func_id           | object | No       | {"value": "","as_default_value": false} value type is string                                               |
| user                 | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| timeout              | object | No       | {"value": null,"as_default_value": false}, value type is number                                            |
| description          | object | No       | {"value": "1","as_default_value": true}}, value type is string                                             |
| bk_start_param_regex | object | No       | {"value": "","as_default_value": false}, value type is string                                              |
| bind_info            | object | No       | {"value":[],"as_default_value": true }, value detailed information see process_property.bind_info.value[n] |

#### Fields that can appear in process_property.bind_info.value[n]

Note:

When modifying bind_info, you must first obtain the content of the original process's bind_info, then modify the
existing bind_info of the process, and pass the modified content to the modification structure.

as_default_value: Whether the value of the process is based on the template value: The value of the process, different
field types have different types

| Name     | Type   | Required | Description                                                                                         |
|----------|--------|----------|-----------------------------------------------------------------------------------------------------|
| enable   | object | No       | {"value": false,"as_default_value": true}, value type is boolean                                    |
| ip       | object | No       | {"value": "1","as_default_value": true}, value type is string                                       |
| port     | object | No       | {"value": "100","as_default_value": true}, value type is string                                     |
| protocol | object | No       | {"value": "1","as_default_value": true}, value type is string                                       |
| row_id   | int    | No       | Unique identifier id, new rows can be set to null, and must maintain the original value for updates |

### Request Example

```python
{
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

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
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

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Updated process template information                                |

#### Explanation of the data field

| Name                | Type   | Description                |
|---------------------|--------|----------------------------|
| id                  | int    | Data ID                    |
| bk_process_name     | string | Process alias              |
| bk_biz_id           | int    | Business ID                |
| service_template_id | int    | Service template ID        |
| property            | object | Properties                 |
| creator             | string | Creator of this data       |
| modifier            | string | Last modifier of this data |
| create_time         | string | Creation time              |
| last_time           | string | Update time                |
| bk_supplier_account | string | Supplier account           |
