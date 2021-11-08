### Functional description

batch create service instance

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| bk_module_id         | int  | Yes   | module ID |
| instances            | array  | Yes   | new service instance data |

#### isntances field description

| Field                |  Type       | Required	   | Description                            |
|---|---|---|---|
|instances.bk_host_id|int|host ID|which host this service instance bind to|
|instances.processes|array|process instance info|process instances in this service instance|


### Request Parameters Example

```json
{
  "bk_biz_id": 1,
  "name": "test1",
  "bk_module_id": 60,
  "instances": [
    {
      "bk_host_id": 2,
      "processes": [
        {
          "process_template_id": 1,
          "process_info": {
            "bk_supplier_account": "0",
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
            "auto_time_gap": 60,
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

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": [53]
}
```

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | new service instance id |

