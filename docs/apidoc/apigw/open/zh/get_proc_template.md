### 描述

获取进程模板信息，url参数中指定进程模板ID

### 输入参数

| 参数名称                | 参数类型 | 必选 | 描述     |
|---------------------|------|----|--------|
| bk_biz_id           | int  | 否  | 业务ID   |
| process_template_id | int  | 是  | 进程模板ID |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "process_template_id": 49
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "id": 49,
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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data 字段说明

| 参数名称                | 参数类型   | 描述          |
|---------------------|--------|-------------|
| id                  | int    | 进程模版id      |
| bk_process_name     | string | 进程别名        |
| bk_biz_id           | int    | 业务id        |
| service_template_id | int    | 服务模版id      |
| property            | object | 进程属性        |
| creator             | string | 本条数据创建者     |
| modifier            | string | 本条数据的最后修改人员 |
| create_time         | string | 创建时间        |
| last_time           | string | 更新时间        |
| bk_supplier_account | string | 开发商账号       |

#### property 字段说明

| 参数名称                 | 参数类型   | 描述      |
|----------------------|--------|---------|
| auto_start           | bool   | 是否自动拉起  |
| bk_biz_id            | int    | 业务id    |
| bk_func_id           | string | 功能ID    |
| bk_func_name         | string | 进程名称    |
| bk_process_id        | int    | 进程id    |
| bk_process_name      | string | 进程别名    |
| bk_start_param_regex | string | 进程启动参数  |
| bk_supplier_account  | string | 开发商账号   |
| create_time          | string | 创建时间    |
| description          | string | 描述      |
| face_stop_cmd        | string | 强制停止命令  |
| last_time            | string | 更新时间    |
| pid_file             | string | PID文件路径 |
| priority             | int    | 启动优先级   |
| proc_num             | int    | 启动数量    |
| reload_cmd           | string | 进程重载命令  |
| restart_cmd          | string | 重启命令    |
| start_cmd            | string | 启动命令    |
| stop_cmd             | string | 停止命令    |
| timeout              | int    | 操作超时时长  |
| user                 | string | 启动用户    |
| work_path            | string | 工作路径    |
| bind_info            | object | 绑定信息    |

#### bind_info 字段说明

| 参数名称     | 参数类型   | 描述                |
|----------|--------|-------------------|
| enable   | bool   | 端口是否启用            |
| ip       | string | 绑定的ip             |
| port     | string | 绑定的端口             |
| protocol | string | 使用的协议             |
| row_id   | int    | 实例化使用的模板行索引，进程内唯一 |
