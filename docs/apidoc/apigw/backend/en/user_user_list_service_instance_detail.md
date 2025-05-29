### Description

Retrieve a list of service instances (with process information) based on the business ID, with optional additional query
conditions such as module ID.

### Parameters

| Name                 | Type   | Required | Description                                                                              |
|----------------------|--------|----------|------------------------------------------------------------------------------------------|
| bk_biz_id            | int    | Yes      | Business ID                                                                              |
| bk_module_id         | int    | No       | Module ID                                                                                |
| bk_host_id           | int    | No       | Host ID (Note: This field is no longer maintained; please use the bk_host_list field)    |
| bk_host_list         | array  | No       | List of host IDs                                                                         |
| service_instance_ids | array  | No       | List of service instance IDs                                                             |
| selectors            | array  | No       | Label filter function; operator optional values: `=`, `!=`, `exists`, `!`, `in`, `notin` |
| page                 | object | Yes      | Pagination parameters                                                                    |

Note: Only one of the parameters `bk_host_list` and `bk_host_id` can be effective, and it is not recommended to
use `bk_host_id` anymore.

#### selectors

| Name     | Type   | Required | Description                                                       |
|----------|--------|----------|-------------------------------------------------------------------|
| key      | string | No       | Field name                                                        |
| operator | string | No       | Operator optional values: `=`, `!=`, `exists`, `!`, `in`, `notin` |
| values   | -      | No       | Different operators correspond to different value formats         |

#### page

| Name  | Type | Required | Description                              |
|-------|------|----------|------------------------------------------|
| start | int  | Yes      | Record start position                    |
| limit | int  | Yes      | Number of records per page, maximum 1000 |

### Request Example

```python
{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 10,
  },
  "bk_module_id": 8,
  "bk_host_list": [11,12],
  "service_instance_ids": [49],
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  }]
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
        "bk_biz_id": 1,
        "id": 49,
        "name": "p1_81",
        "service_template_id": 50,
        "bk_host_id": 11,
        "bk_module_id": 56,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-07-22T09:54:50.906+08:00",
        "last_time": "2019-07-22T09:54:50.906+08:00",
        "bk_supplier_account": "0",
        "service_category_id": 22,
        "process_instances": [
          {
            "process": {
              "proc_num": 0,
              "stop_cmd": "",
              "restart_cmd": "",
              "face_stop_cmd": "",
              "bk_process_id": 43,
              "bk_func_name": "p1",
              "work_path": "",
              "priority": 0,
              "reload_cmd": "",
              "bk_process_name": "p1",
              "pid_file": "",
              "auto_start": false,
              "last_time": "2019-07-22T09:54:50.927+08:00",
              "create_time": "2019-07-22T09:54:50.927+08:00",
              "bk_biz_id": 3,
              "start_cmd": "",
              "user": "",
              "timeout": 0,
              "description": "",
              "bk_supplier_account": "0",
              "bk_start_param_regex": "",
              "bind_info": [
                {
                    "enable": true,
                    "ip": "127.0.0.1",
                    "port": "80",
                    "protocol": "1",
                    "template_row_id": 1234
                }
              ]
            },
            "relation": {
              "bk_biz_id": 1,
              "bk_process_id": 43,
              "service_instance_id": 49,
              "process_template_id": 48,
              "bk_host_id": 11,
              "bk_supplier_account": "0"
            }
          }
        ]
      }
    ]
  }
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

| Name  | Type  | Description              |
|-------|-------|--------------------------|
| count | int   | Total number of records  |
| info  | array | List of returned results |

#### info Field Explanation

| Name                       | Type   | Description                                                                            |
|----------------------------|--------|----------------------------------------------------------------------------------------|
| id                         | int    | Service instance ID                                                                    |
| name                       | array  | Service instance name                                                                  |
| service_template_id        | int    | Service template ID                                                                    |
| bk_host_id                 | int    | Host ID                                                                                |
| bk_host_innerip            | string | Host IP                                                                                |
| bk_module_id               | int    | Module ID                                                                              |
| creator                    | string | Creator of this data                                                                   |
| modifier                   | string | Last modifier of this data                                                             |
| create_time                | string | Creation time                                                                          |
| last_time                  | string | Update time                                                                            |
| bk_supplier_account        | string | Supplier account                                                                       |
| service_category_id        | int    | Service category ID                                                                    |
| process_instances          | array  | Process instance information                                                           |
| bk_biz_id                  | int    | Business ID                                                                            |
| process_instances.process  | object | Process instance details, process attribute fields                                     |
| process_instances.relation | object | Relationship information of the process instance, such as host ID, process template ID |

#### data.info.process_instances[x].process Field Explanation

| Name                 | Type   | Description                    |
|----------------------|--------|--------------------------------|
| auto_start           | bool   | Whether to start automatically |
| auto_time_gap        | int    | Startup interval               |
| bk_biz_id            | int    | Business ID                    |
| bk_func_id           | string | Function ID                    |
| bk_func_name         | string | Process name                   |
| bk_process_id        | int    | Process ID                     |
| bk_process_name      | string | Process alias                  |
| bk_start_param_regex | string | Process startup parameters     |
| bk_supplier_account  | string | Supplier account               |
| create_time          | string | Creation time                  |
| description          | string | Description                    |
| face_stop_cmd        | string | Force stop command             |
| last_time            | string | Update time                    |
| pid_file             | string | PID file path                  |
| priority             | int    | Startup priority               |
| proc_num             | int    | Number of startups             |
| reload_cmd           | string | Process reload command         |
| restart_cmd          | string | Restart command                |
| start_cmd            | string | Start command                  |
| stop_cmd             | string | Stop command                   |
| timeout              | int    | Operation timeout duration     |
| user                 | string | Start user                     |
| work_path            | string | Working directory              |
| bind_info            | object | Binding information            |

#### data.info.process_instances[x].process.bind_info Field Explanation

| Name            | Type   | Description                                   |
|-----------------|--------|-----------------------------------------------|
| enable          | bool   | Whether the port is enabled                   |
| ip              | string | Bound IP                                      |
| port            | string | Bound port                                    |
| protocol        | string | Used protocol                                 |
| template_row_id | int    | Template row index, unique within the process |

#### data.info.process_instances[x].relation Field Explanation

| Name                | Type   | Description         |
|---------------------|--------|---------------------|
| bk_biz_id           | int    | Business ID         |
| bk_process_id       | int    | Process ID          |
| service_instance_id | int    | Service instance ID |
| process_template_id | int    | Process template ID |
| bk_host_id          | int    | Host ID             |
| bk_supplier_account | string | Supplier account    |
