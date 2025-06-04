### 描述

查询某业务下进程ID对应的进程详情 (v3.9.8)

### 输入参数

| 参数名称           | 参数类型  | 必选 | 描述                                                                              |
|----------------|-------|----|---------------------------------------------------------------------------------|
| bk_biz_id      | int   | 是  | 进程所在的业务ID                                                                       |
| bk_process_ids | array | 是  | 进程ID列表，最多传500个                                                                  |
| fields         | array | 否  | 进程属性列表，控制返回结果的进程实例信息里有哪些字段，能够加速接口请求和减少网络流量传输<br>为空时返回进程所有字段,bk_process_id为必返回字段 |

### 调用示例

```json
{
    "bk_biz_id":1,
    "bk_process_ids": [
        43,
        44
    ],
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_name"
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
        "data": [
        {
            "auto_start": null,
            "bind_info": [
                {
                    "enable": true,
                    "ip": "127.0.0.1",
                    "port": "9898",
                    "protocol": "1",
                    "template_row_id": 1
                }
            ],
            "bk_biz_id": 3,
            "bk_created_at": "2023-11-15T10:37:39.384+08:00",
            "bk_created_by": "admin",
            "bk_func_name": "delete",
            "bk_process_id": 57,
            "bk_process_name": "delete-aa",
            "bk_start_check_secs": null,
            "bk_start_param_regex": "",
            "bk_supplier_account": "0",
            "bk_updated_at": "2023-11-15T17:19:18.1+08:00",
            "bk_updated_by": "admin",
            "create_time": "2023-11-15T10:37:39.384+08:00",
            "description": "",
            "face_stop_cmd": "",
            "last_time": "2023-11-15T17:19:18.1+08:00",
            "pid_file": "",
            "priority": null,
            "proc_num": null,
            "reload_cmd": "",
            "restart_cmd": "",
            "service_instance_id": 57,
            "start_cmd": "",
            "stop_cmd": "",
            "timeout": null,
            "user": "",
            "work_path": ""
        }
    ],
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | array  | 请求返回的数据                    |

#### data

| 参数名称                 | 参数类型   | 描述      |
|----------------------|--------|---------|
| auto_start           | bool   | 是否自动拉起  |
| bk_biz_id            | int    | 业务id    |
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
| bk_created_at        | string | 创建时间    |
| bk_created_by        | string | 创建人     |
| bk_updated_at        | string | 更新时间    |
| bk_updated_by        | string | 更新人     |
| bind_info            | object | 绑定信息    |

#### data.info.property.bind_info.value 字段说明

| 参数名称     | 参数类型   | 描述          |
|----------|--------|-------------|
| enable   | object | 端口是否启用      |
| ip       | object | 绑定的ip       |
| port     | object | 绑定的端口       |
| protocol | object | 使用的协议       |
| row_id   | int    | 模板行索引，进程内唯一 |
