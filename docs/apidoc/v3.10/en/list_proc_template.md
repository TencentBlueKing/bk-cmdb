### Functional description

list process templates

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_biz_id | int        | Yes     | business id |
| service_template_id | int  | No   | service template id, must have at least one of service_template_id or process_template_ids |
| process_template_ids | int array  | No   | the array of process template id , must have at least one of service_template_id or process_template_ids|

### Request Parameters Example

```json
{
    "bk_biz_id": 1,
    "service_template_id": 51,
    "process_template_ids": [
        50
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "count": 1,
        "info": [
            {
                "id": 50,
                "bk_process_name": "p1",
                "bk_biz_id": 1,
                "service_template_id": 51,
                "property": {
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
                },
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2019-06-19T15:24:04.763+08:00",
                "last_time": "2019-06-19T15:24:04.763+08:00",
                "bk_supplier_account": "0"
            }
        ]
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

#### Data field description

| Field       | Type     | Description         |
|---|---|---|
|count|integer|total count|
|info|array|response data|

#### Info field description

| Field       | Type     | Description         |
|---|---|---|
|id|integer| process template ID|
|bk_process_name|string|process template name|
|property|object|process template properties |



#### data.info[x].property.bind_info.value[n] description
| Field|Type|Description|
|---|---|---|
|enable|object|Whether the port is enabled|
|ip|object|bind ip|
|port|object|bind port|
|protocol|object|protocol used|
|row_id|int|template row index, unique in process|

