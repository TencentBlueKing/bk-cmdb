### 功能描述

批量创建进程模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id  | int     |是     | 业务ID       |
| service_template_id            | int  | 否   | 服务模板ID |
| processes         | array  | 是   | 进程模板信息 |


### 请求参数示例

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

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [[52]]
}
```

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | array | 成功创建的进程模板ID |
