### 描述

批量创建服务实例，如果模块绑定了服务模板，则服务实例也会根据模板创建，创建服务实例的进程参数内也必须提供每个进程对应的进程模板ID(
权限：服务实例新建权限)

### 输入参数

| 参数名称         | 参数类型  | 必选 | 描述                  |
|--------------|-------|----|---------------------|
| bk_module_id | int   | 是  | 模块ID                |
| instances    | array | 是  | 需要创建的服务实例信息,最大值为100 |
| bk_biz_id    | int   | 是  | 业务ID                |

#### instances 字段说明

| 参数名称                                    | 参数类型   | 必选 | 描述                                                                  |
|-----------------------------------------|--------|----|---------------------------------------------------------------------|
| instances.bk_host_id                    | int    | 是  | 主机ID,服务实例绑定的主机ID                                                    |
| instances.service_instance_name         | string | 否  | 服务实例名称，不填则会使用主机IP加进程名称加服务绑定的端口作为名称，如“123.123.123.123_job_java_80”形式 |
| instances.processes                     | array  | 是  | 进程信息,服务实例下新建的进程信息                                                   |
| instances.processes.process_template_id | int    | 是  | 进程模板ID,如果模块没有绑定服务模板则填0                                              |
| instances.processes.process_info        | object | 是  | 进程实例信息,如果进程绑定有模板，则仅模板中没有锁定的字段有效                                     |

#### processes 字段说明

| 参数名称                | 参数类型   | 必选 | 描述      |
|---------------------|--------|----|---------|
| process_template_id | int    | 是  | 进程模版id  |
| auto_start          | bool   | 否  | 是否自动拉起  |
| bk_biz_id           | int    | 否  | 业务id    |
| bk_func_id          | string | 否  | 功能ID    |
| bk_func_name        | string | 否  | 进程名称    |
| bk_process_id       | int    | 否  | 进程id    |
| bk_process_name     | string | 否  | 进程别名    |
| bk_supplier_account | string | 否  | 开发商账号   |
| face_stop_cmd       | string | 否  | 强制停止命令  |
| pid_file            | string | 否  | PID文件路径 |
| priority            | int    | 否  | 启动优先级   |
| proc_num            | int    | 否  | 启动数量    |
| reload_cmd          | string | 否  | 进程重载命令  |
| restart_cmd         | string | 否  | 重启命令    |
| start_cmd           | string | 否  | 启动命令    |
| stop_cmd            | string | 否  | 停止命令    |
| timeout             | int    | 否  | 操作超时时长  |
| user                | string | 否  | 启动用户    |
| work_path           | string | 否  | 工作路径    |
| process_info        | object | 是  | 进程信息    |

#### process_info 字段说明

| 参数名称                | 参数类型   | 必选 | 描述    |
|---------------------|--------|----|-------|
| bind_info           | object | 是  | 绑定信息  |
| bk_supplier_account | string | 是  | 开发商账号 |

#### bind_info 字段说明

| 参数名称            | 参数类型   | 必选 | 描述                |
|-----------------|--------|----|-------------------|
| enable          | bool   | 是  | 端口是否启用            |
| ip              | string | 是  | 绑定的ip             |
| port            | string | 是  | 绑定的端口             |
| protocol        | string | 是  | 使用的协议             |
| template_row_id | int    | 是  | 实例化使用的模板行索引，进程内唯一 |

### 调用示例

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

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [53]
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 新建的服务实例ID列表                |
