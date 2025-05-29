### 描述

根据服务模板ID查询进程模板信息

### 输入参数

| 参数名称                 | 参数类型  | 必选 | 描述                                                            |
|----------------------|-------|----|---------------------------------------------------------------|
| bk_biz_id            | int   | 是  | 业务id                                                          |
| service_template_id  | int   | 否  | 服务模板ID，service_template_id和process_template_ids至少传一个          |
| process_template_ids | array | 否  | 进程模板ID数组，最多200个，service_template_id和process_template_ids至少传一个 |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "service_template_id": 51,
    "process_template_ids": [
        50
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
  "data": {
    "count": 1,
    "info": [
      {
        "id": 6,
        "bk_process_name": "red-1",
        "bk_biz_id": 3,
        "service_template_id": 5,
        "property": {
          "proc_num": {
            "value": null,
            "as_default_value": true
          },
          "stop_cmd": {
            "value": "",
            "as_default_value": true
          },
          "restart_cmd": {
            "value": "",
            "as_default_value": true
          },
          "face_stop_cmd": {
            "value": "",
            "as_default_value": true
          },
          "bk_func_name": {
            "value": "red-1",
            "as_default_value": true
          },
          "work_path": {
            "value": "",
            "as_default_value": true
          },
          "priority": {
            "value": null,
            "as_default_value": true
          },
          "reload_cmd": {
            "value": "",
            "as_default_value": true
          },
          "bk_process_name": {
            "value": "red-1",
            "as_default_value": true
          },
          "pid_file": {
            "value": "",
            "as_default_value": true
          },
          "auto_start": {
            "value": null,
            "as_default_value": null
          },
          "bk_start_check_secs": {
            "value": null,
            "as_default_value": true
          },
          "start_cmd": {
            "value": "",
            "as_default_value": true
          },
          "user": {
            "value": "",
            "as_default_value": true
          },
          "timeout": {
            "value": null,
            "as_default_value": true
          },
          "description": {
            "value": "",
            "as_default_value": true
          },
          "bk_start_param_regex": {
            "value": "",
            "as_default_value": true
          },
          "bind_info": {
            "value": [
              {
                "enable": {
                  "value": true,
                  "as_default_value": true
                },
                "ip": {
                  "value": "1",
                  "as_default_value": true
                },
                "port": {
                  "value": "9583",
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
        "create_time": "2023-11-15T02:10:04.619Z",
        "last_time": "2023-11-15T02:10:04.619Z",
        "bk_supplier_account": "0"
      }
    ]
  },
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

| 参数名称  | 参数类型  | 描述   |
|-------|-------|------|
| count | int   | 总数   |
| info  | array | 返回结果 |

#### info 字段说明

| 参数名称                | 参数类型   | 描述          |
|---------------------|--------|-------------|
| id                  | int    | 进程模板ID      |
| bk_process_name     | string | 进程模板名称      |
| property            | object | 进程模板属性      |
| bk_biz_id           | int    | 业务ID        |
| service_template_id | int    | 服务模版ID      |
| creator             | string | 本条数据创建者     |
| modifier            | string | 本条数据的最后修改人员 |
| create_time         | string | 创建时间        |
| last_time           | string | 更新时间        |
| bk_supplier_account | string | 开发商账号       |

#### data.info[x].property

as_default_value 进程的值是否以模板为准

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
| bind_info            | object | 绑定信息    |

#### data.info[x].property.bind_info.value[n] 字段说明

| 参数名称     | 参数类型   | 描述          |
|----------|--------|-------------|
| enable   | object | 端口是否启用      |
| ip       | object | 绑定的ip       |
| port     | object | 绑定的端口       |
| protocol | object | 使用的协议       |
| row_id   | int    | 模板行索引，进程内唯一 |
