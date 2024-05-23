### Description

Batch create service instances. If the module is bound to a service template, the service instances will also be created
based on the template. The process parameters for creating service instances must also provide the process template ID
corresponding to each process (Permission: Service Instance Creation Permission).

### Parameters

| Name         | Type  | Required | Description                                                                 |
|--------------|-------|----------|-----------------------------------------------------------------------------|
| bk_module_id | int   | Yes      | Module ID                                                                   |
| instances    | array | Yes      | Information of service instances to be created, with a maximum value of 100 |
| bk_biz_id    | int   | Yes      | Business ID                                                                 |

#### Explanation of the instances Field

| Name                                    | Type   | Required | Description                                                                                                                                                                  |
|-----------------------------------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| instances.bk_host_id                    | int    | Yes      | Host ID, the host ID bound to the service instance                                                                                                                           |
| instances.service_instance_name         | string | No       | Service instance name. If not filled, the host IP plus the process name plus the service binding port will be used as the name, in the form of "123.123.123.123_job_java_80" |
| instances.processes                     | array  | Yes      | Process information, information of processes newly created under the service instance                                                                                       |
| instances.processes.process_template_id | int    | Yes      | Process template ID. If the module is not bound to a service template, fill in 0                                                                                             |
| instances.processes.process_info        | object | Yes      | Process instance information. If the process is bound to a template, only the fields not locked in the template are valid                                                    |

#### Explanation of the processes Field

| Name                | Type   | Required | Description                    |
|---------------------|--------|----------|--------------------------------|
| process_template_id | int    | Yes      | Process template ID            |
| auto_start          | bool   | No       | Whether to start automatically |
| bk_biz_id           | int    | No       | Business ID                    |
| bk_func_id          | string | No       | Function ID                    |
| bk_func_name        | string | No       | Process name                   |
| bk_process_id       | int    | No       | Process ID                     |
| bk_process_name     | string | No       | Process alias                  |
| bk_supplier_account | string | No       | Vendor account                 |
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
| process_info        | object | Yes      | Process information            |

#### Explanation of the process_info Field

| Name                | Type   | Required | Description         |
|---------------------|--------|----------|---------------------|
| bind_info           | object | Yes      | Binding information |
| bk_supplier_account | string | Yes      | Vendor account      |

#### Explanation of the bind_info Field

| Name            | Type   | Required | Description                                                          |
|-----------------|--------|----------|----------------------------------------------------------------------|
| enable          | bool   | Yes      | Whether the port is enabled                                          |
| ip              | string | Yes      | Bound IP                                                             |
| port            | string | Yes      | Bound port                                                           |
| protocol        | string | Yes      | Used protocol                                                        |
| template_row_id | int    | Yes      | Template row index used for instantiation, unique within the process |

### Request Example

```json
{
  "bk_biz_id": 1,
  "bk_module_id": 60,
  "instances": [
    {
      "bk_host_id": 2,
      "service_instance_name": "test",
      "processes": [
        {
          "process_template_id": 1,
          "process_info": {
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

### Response Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [53]
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | List of newly created service instance IDs                                  |
