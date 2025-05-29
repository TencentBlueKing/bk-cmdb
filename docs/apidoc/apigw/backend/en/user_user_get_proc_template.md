### Description

Get process template information by specifying the process template ID in the URL.

### Parameters

| Name                | Type | Required | Description         |
|---------------------|------|----------|---------------------|
| bk_biz_id           | int  | No       | Business ID         |
| process_template_id | int  | Yes      | Process template ID |

### Request Example

```python
{
  "bk_biz_id": 1,
  "process_template_id": 49
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

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error    |
| message    | string | Error message returned in case of request failure                |
| permission | object | Permission information                                           |
| data       | object | Data returned by the request                                     |

#### data Field Description

| Name                | Type   | Description                |
|---------------------|--------|----------------------------|
| id                  | int    | Process template ID        |
| bk_process_name     | string | Process alias              |
| bk_biz_id           | int    | Business ID                |
| service_template_id | int    | Service template ID        |
| property            | object | Process properties         |
| creator             | string | Creator of this data       |
| modifier            | string | Last modifier of this data |
| create_time         | string | Creation time              |
| last_time           | string | Last update time           |
| bk_supplier_account | string | Supplier account           |

#### property Field Description

| Name                 | Type   | Description                  |
|----------------------|--------|------------------------------|
| auto_start           | bool   | Auto start flag              |
| bk_biz_id            | int    | Business ID                  |
| bk_func_id           | string | Function ID                  |
| bk_func_name         | string | Process name                 |
| bk_process_id        | int    | Process ID                   |
| bk_process_name      | string | Process alias                |
| bk_start_param_regex | string | Process start parameters     |
| bk_supplier_account  | string | Supplier account             |
| create_time          | string | Creation time                |
| description          | string | Description                  |
| face_stop_cmd        | string | Force stop command           |
| last_time            | string | Last update time             |
| pid_file             | string | PID file path                |
| priority             | int    | Start priority               |
| proc_num             | int    | Number of instances to start |
| reload_cmd           | string | Process reload command       |
| restart_cmd          | string | Restart command              |
| start_cmd            | string | Start command                |
| stop_cmd             | string | Stop command                 |
| timeout              | int    | Operation timeout            |
| user                 | string | Start user                   |
| work_path            | string | Working directory            |
| bind_info            | object | Binding information          |

#### bind_info Field Description

| Name     | Type   | Description                                                          |
|----------|--------|----------------------------------------------------------------------|
| enable   | bool   | Whether the port is enabled                                          |
| ip       | string | Bound IP                                                             |
| port     | string | Bound port                                                           |
| protocol | string | Used protocol                                                        |
| row_id   | int    | Template row index used for instantiation, unique within the process |
