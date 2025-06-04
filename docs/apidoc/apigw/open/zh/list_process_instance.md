### 描述

根据服务实例ID查询进程实例列表

### 输入参数

| 参数名称                | 参数类型 | 必选 | 描述     |
|---------------------|------|----|--------|
| bk_biz_id           | int  | 是  | 业务ID   |
| service_instance_id | int  | 是  | 服务实例ID |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "service_instance_id": 54
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
      "property": {
        "auto_start": false,
        "auto_time_gap": 60,
        "bk_biz_id": 2,
        "bk_func_id": "",
        "bk_func_name": "java",
        "bk_process_id": 46,
        "bk_process_name": "job_java",
        "bk_start_param_regex": "",
        "bk_supplier_account": "0",
        "create_time": "2019-06-05T14:59:12.065+08:00",
        "description": "",
        "face_stop_cmd": "",
        "last_time": "2019-06-05T14:59:12.065+08:00",
        "pid_file": "",
        "priority": 1,
        "proc_num": 1,
        "reload_cmd": "",
        "restart_cmd": "",
        "start_cmd": "",
        "stop_cmd": "",
        "timeout": 30,
        "user": "",
        "work_path": "/data/bkee",
        "bind_info": [
            {
                "enable": false,  
                "ip": "127.0.0.1",  
                "port": "100",  
                "protocol": "1", 
                "template_row_id": 1  
            }
        ]
      },
      "relation": {
        "bk_biz_id": 1,
        "bk_process_id": 46,
        "service_instance_id": 54,
        "process_template_id": 1,
        "bk_host_id": 1,
        "bk_supplier_account": ""
      }
    }
  ]
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

#### data 字段说明

| 参数名称     | 参数类型   | 描述           |
|----------|--------|--------------|
| property | object | 进程属性信息       |
| relation | object | 进程与服务实例的关联信息 |

#### data[x].property 字段说明

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

#### data[x].property.bind_info[n] 字段说明

| 参数名称            | 参数类型   | 描述                |
|-----------------|--------|-------------------|
| enable          | bool   | 端口是否启用            |
| ip              | string | 绑定的ip             |
| port            | string | 绑定的端口             |
| protocol        | string | 使用的协议             |
| template_row_id | int    | 实例化使用的模板行索引，进程内唯一 |

#### data[x].relation 字段说明

| 参数名称                | 参数类型   | 描述     |
|---------------------|--------|--------|
| bk_biz_id           | int    | 业务id   |
| bk_process_id       | int    | 进程id   |
| service_instance_id | int    | 服务实例id |
| process_template_id | int    | 进程模版id |
| bk_host_id          | int    | 主机id   |
| bk_supplier_account | string | 开发商账号  |
