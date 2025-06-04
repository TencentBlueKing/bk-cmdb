### 描述

根据服务实例ID和进程实例属性值创建进程实例(权限：服务实例编辑权限)

### 输入参数

| 参数名称                | 参数类型  | 必选 | 描述                     |
|---------------------|-------|----|------------------------|
| service_instance_id | int   | 是  | 服务实例ID                 |
| processes           | array | 是  | 进程实例中不同于模板的属性值,最大值为100 |

#### process_info 字段说明

| 参数名称                | 参数类型   | 必选 | 描述      |
|---------------------|--------|----|---------|
| auto_start          | bool   | 否  | 是否自动拉起  |
| auto_time_gap       | int    | 否  | 拉起间隔    |
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
| bind_info           | object | 否  | 绑定信息    |

#### bind_info 字段说明

| 参数名称     | 参数类型   | 必选 | 描述                |
|----------|--------|----|-------------------|
| enable   | bool   | 否  | 端口是否启用            |
| ip       | string | 否  | 绑定的ip             |
| port     | string | 否  | 绑定的端口             |
| protocol | string | 否  | 使用的协议             |
| row_id   | int    | 否  | 实例化使用的模板行索引，进程内唯一 |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "service_instance_id": 48,
  "processes": [
    {
      "process_info": {
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
        "bk_process_name": "job_java",
        "user": "",
        "proc_num": 1,
        "priority": 1,
        "bk_biz_id": 2,
        "bk_start_param_regex": "",
        "bk_process_id": 1,
        "bind_info": [
          {
              "enable": false,
              "ip": "127.0.0.1",
              "port": "80",
              "protocol": "1",
              "template_row_id": 1234
          }
        ]
      }
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
  "data": [64]
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 新建的进程实例ID列表                |
