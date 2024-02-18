### Description

Batch Create Process Templates (Permission: Service Template Editing Permission)

### Parameters

| Name                | Type  | Required | Description                                               |
|---------------------|-------|----------|-----------------------------------------------------------|
| bk_biz_id           | int   | Yes      | Business ID                                               |
| service_template_id | int   | No       | Service template ID                                       |
| processes           | array | Yes      | Process template information, with a maximum value of 100 |

#### processes

as_default_value: Whether the value of the process is based on the template

| Name                | Type   | Required | Description                    |
|---------------------|--------|----------|--------------------------------|
| auto_start          | bool   | No       | Whether to start automatically |
| bk_biz_id           | int    | No       | Business ID                    |
| bk_func_id          | string | No       | Function ID                    |
| bk_func_name        | string | No       | Process name                   |
| bk_process_id       | int    | No       | Process ID                     |
| bk_process_name     | string | No       | Process alias                  |
| bk_supplier_account | string | No       | Supplier account               |
| face_stop_cmd       | string | No       | Force stop command             |
| pid_file            | string | No       | PID file path                  |
| priority            | int    | No       | Startup priority               |
| proc_num            | int    | No       | Number of startups             |
| reload_cmd          | string | No       | Process reload command         |
| restart_cmd         | string | No       | Restart command                |
| start_cmd           | string | No       | Startup command                |
| stop_cmd            | string | No       | Stop command                   |
| timeout             | int    | No       | Operation timeout duration     |
| user                | string | No       | Startup user                   |
| work_path           | string | No       | Working directory              |
| bind_info           | object | No       | Binding information            |

#### bind_info Field Description

| Name     | Type   | Required | Description                                                          |
|----------|--------|----------|----------------------------------------------------------------------|
| enable   | bool   | No       | Whether the port is enabled                                          |
| ip       | string | No       | Bound IP                                                             |
| port     | string | No       | Bound port                                                           |
| protocol | string | No       | Protocol used                                                        |
| row_id   | int    | No       | Template row index used for instantiation, unique within the process |

### Request Example

```json
{ 
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

### Response Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [52]
}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | array  | IDs of successfully created process templates                     |
