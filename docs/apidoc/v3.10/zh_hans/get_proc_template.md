### 功能描述

获取单个进程模板信息，url参数中指定进程模板ID

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|  
| bk_biz_id            | int        |是     | 业务ID        |
| process_template_id  | int        |是     | 进程模板ID     |

### 请求参数示例

```python
{
  "bk_biz_id": 1,
  "process_template_id": 49,
  "bk_supplier_account": "0"
}
```

### 返回结果示例

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

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | object | 请求返回的数据 |
