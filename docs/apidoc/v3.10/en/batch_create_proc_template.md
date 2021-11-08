### Functional description

batch create process templates

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| bk_biz_id  | int     |Yes     | Business ID       |
| service_template_id            | int  | No   | Service Template ID |
| processes         | array  | Yes   | Process Template Info |


### Request Parameters Example

```json
{
  "bk_biz_id": 1,
  "service_template_id": 1,
  "processes": [
    {
      "spec": {
          "auto_start": {
            "as_default_value": true,
            "value": false
          },
          "auto_time_gap": {
            "as_default_value": false
          },
          "bind_info": {
            "as_default_value": true,
            "value": [
              {
                "ip": {
                  "value": "2",
                  "as_default_value": true
                },
                "port": {
                  "value": "1",
                  "as_default_value": false
                },
                "protocol": {
                  "value": "2",
                  "as_default_value": false
                },
                "enable": {
                  "value": false,
                  "as_default_value": false
                }
              }
            ]
          },
          "bk_biz_id": {
            "as_default_value": true,
            "value": 2
          },
          "bk_func_id": {
            "as_default_value": true,
            "value": "1"
          },
          "bk_func_name": {
            "as_default_value": true,
            "value": "nginx"
          },
          "bk_process_id": {
            "as_default_value": true,
            "value": 3
          },
          "bk_process_name": {
            "as_default_value": true,
            "value": ""
          },
          "bk_supplier_account": {
            "as_default_value": true,
            "value": ""
          },
          "create_time": {
            "as_default_value": true,
            "value": "2019-05-06T07:12:35.082Z"
          },
          "description": {
            "as_default_value": true,
            "value": "a simple description"
          },
          "face_stop_cmd": {
            "as_default_value": true,
            "value": "./stop.sh"
          },
          "last_time": {
            "as_default_value": true,
            "value": "2019-05-06T07:12:35.082Z"
          },
          "pid_file": {
            "as_default_value": true,
            "value": ""
          },
          "priority": {
            "as_default_value": true,
            "value": 1
          },
          "proc_num": {
            "as_default_value": true,
            "value": 1
          },
          "reload_cmd": {
            "as_default_value": true,
            "value": ""
          },
          "restart_cmd": {
            "as_default_value": true,
            "value": "./restart.sh"
          },
          "start_cmd": {
            "as_default_value": true,
            "value": "./start.sh"
          },
          "stop_cmd": {
            "as_default_value": true,
            "value": "./stop.sh"
          },
          "timeout": {
            "as_default_value": true,
            "value": 60
          },
          "user": {
            "as_default_value": true,
            "value": ""
          },
          "work_path": {
            "as_default_value": true,
            "value": "/data/bkee"
          },
          "bk_start_param_regex": {
            "as_default_value": true,
            "value": ""
          }
        }
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
  "permission": null,
  "data": [[52]]
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
