### Description

Query process template information based on service template ID

### Parameters

| Name                 | Type  | Required | Description                                                                                                           |
|----------------------|-------|----------|-----------------------------------------------------------------------------------------------------------------------|
| bk_biz_id            | int   | Yes      | Business ID                                                                                                           |
| service_template_id  | int   | No       | Service template ID, at least one of service_template_id and process_template_ids must be passed                      |
| process_template_ids | array | No       | Array of process template IDs, up to 200, at least one of service_template_id and process_template_ids must be passed |

### Request Example

```json
{
    "bk_biz_id": 1,
    "service_template_id": 51,
    "process_template_ids": [
        50
    ]
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
    "count": 1,
    "info": [
      {
        "id": 6,
        "bk_process_name": "red-1",
        "bk_biz_id": 3,
        "service_template_id": 5,
        "property": {
          "proc_num": {
            "value": null,
            "as_default_value": true
          },
          "stop_cmd": {
            "value": "",
            "as_default_value": true
          },
          "restart_cmd": {
            "value": "",
            "as_default_value": true
          },
          "face_stop_cmd": {
            "value": "",
            "as_default_value": true
          },
          "bk_func_name": {
            "value": "red-1",
            "as_default_value": true
          },
          "work_path": {
            "value": "",
            "as_default_value": true
          },
          "priority": {
            "value": null,
            "as_default_value": true
          },
          "reload_cmd": {
            "value": "",
            "as_default_value": true
          },
          "bk_process_name": {
            "value": "red-1",
            "as_default_value": true
          },
          "pid_file": {
            "value": "",
            "as_default_value": true
          },
          "auto_start": {
            "value": null,
            "as_default_value": null
          },
          "bk_start_check_secs": {
            "value": null,
            "as_default_value": true
          },
          "start_cmd": {
            "value": "",
            "as_default_value": true
          },
          "user": {
            "value": "",
            "as_default_value": true
          },
          "timeout": {
            "value": null,
            "as_default_value": true
          },
          "description": {
            "value": "",
            "as_default_value": true
          },
          "bk_start_param_regex": {
            "value": "",
            "as_default_value": true
          },
          "bind_info": {
            "value": [
              {
                "enable": {
                  "value": true,
                  "as_default_value": true
                },
                "ip": {
                  "value": "1",
                  "as_default_value": true
                },
                "port": {
                  "value": "9583",
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
        "create_time": "2023-11-15T02:10:04.619Z",
        "last_time": "2023-11-15T02:10:04.619Z",
        "bk_supplier_account": "0"
      }
    ]
  },
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data Field Explanation

| Name  | Type  | Description         |
|-------|-------|---------------------|
| count | int   | Number of records   |
| info  | array | Result of the query |

#### info Field Explanation

| Name                | Type   | Description                 |
|---------------------|--------|-----------------------------|
| id                  | int    | Process template ID         |
| bk_process_name     | string | Process template name       |
| property            | object | Process template properties |
| bk_biz_id           | int    | Business ID                 |
| service_template_id | int    | Service template ID         |
| creator             | string | Creator of this data        |
| modifier            | string | Last modifier of this data  |
| create_time         | string | Creation time               |
| last_time           | string | Update time                 |
| bk_supplier_account | string | Supplier account            |

#### data.info[x].property

as_default_value Whether the value of the process is based on the template

| Name                 | Type   | Description                    |
|----------------------|--------|--------------------------------|
| auto_start           | bool   | Whether to start automatically |
| bk_biz_id            | int    | Business ID                    |
| bk_func_name         | string | Process name                   |
| bk_process_id        | int    | Process ID                     |
| bk_process_name      | string | Process alias                  |
| bk_start_param_regex | string | Process start parameters       |
| bk_supplier_account  | string | Supplier account               |
| create_time          | string | Creation time                  |
| description          | string | Description                    |
| face_stop_cmd        | string | Forced stop command            |
| last_time            | string | Update time                    |
| pid_file             | string | PID file path                  |
| priority             | int    | Start priority                 |
| proc_num             | int    | Number of starts               |
| reload_cmd           | string | Process reload command         |
| restart_cmd          | string | Restart command                |
| start_cmd            | string | Start command                  |
| stop_cmd             | string | Stop command                   |
| timeout              | int    | Operation timeout duration     |
| user                 | string | Start user                     |
| work_path            | string | Working path                   |
| bind_info            | object | Binding information            |

#### data.info[x].property.bind_info.value[n] Field Explanation

| Name     | Type   | Description                                   |
|----------|--------|-----------------------------------------------|
| enable   | object | Whether the port is enabled                   |
| ip       | object | Bound IP                                      |
| port     | object | Bound port                                    |
| protocol | object | Used protocol                                 |
| row_id   | int    | Template row index, unique within the process |
