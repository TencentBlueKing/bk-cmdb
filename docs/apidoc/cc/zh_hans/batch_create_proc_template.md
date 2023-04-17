### 功能描述

批量创建进程模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id  | int     |是     | 业务ID       |
| service_template_id            | int  | 否   | 服务模板ID |
| processes         | array  | 是   | 进程模板信息,最大值为100 |


#### processes 
as_default_value: 进程的值是否以模板为准

| 字段|类型|必选|说明|
|---|---|---|---|
|auto_start|bool|否|是否自动拉起|
|auto_time_gap|int|否|拉起间隔|
|bk_biz_id|int|否|业务id|
|bk_func_id|string|否|功能ID|
|bk_func_name|string|否|进程名称|
|bk_process_id|int|否|进程id|
|bk_process_name|string|否|进程别名|_``_
|bk_supplier_account|string|否|开发商账号|
|face_stop_cmd|string|否|强制停止命令|
|pid_file|string|否|PID文件路径|
|priority|int|否|启动优先级|
|proc_num|int|否|启动数量|
|reload_cmd|string|否|进程重载命令|
|restart_cmd|string|否|重启命令|
|start_cmd|string|否|启动命令|
|stop_cmd|string|否|停止命令|
|timeout|int|否|操作超时时长|
|user|string|否|启动用户|
|work_path|string|否|工作路径|
|bind_info|object|否|绑定信息|

#### bind_info 字段说明
| 字段|类型|必选|说明|
|---|---|---|---|
|enable|bool|否|端口是否启用|
|ip|string|否|绑定的ip|
|port|string|否|绑定的端口|
|protocol|string|否|使用的协议|
|row_id|int|否|实例化使用的模板行索引，进程内唯一|

### 请求参数示例

```json
{ 
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
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
  "request_id": "e43da4ef221746868dc4c837d36f3807",
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
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | array | 成功创建的进程模板ID |
