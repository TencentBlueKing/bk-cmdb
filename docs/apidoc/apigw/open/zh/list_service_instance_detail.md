### 描述

根据业务id查询服务实例列表(带进程信息)，可再附加上模块id等查询条件

### 输入参数

| 参数名称                 | 参数类型   | 必选 | 描述                                                        |
|----------------------|--------|----|-----------------------------------------------------------|
| bk_biz_id            | int    | 是  | 业务ID                                                      |
| bk_module_id         | int    | 否  | 模块ID                                                      |
| bk_host_id           | int    | 否  | 主机ID, 注意：该字段不再维护，请用bk_host_list字段                         |
| bk_host_list         | array  | 否  | 主机ID列表                                                    |
| service_instance_ids | int    | 否  | 服务实例ID列表                                                  |
| selectors            | int    | 否  | label过滤功能，operator可选值: `=`,`!=`,`exists`,`!`,`in`,`notin` |
| page                 | object | 是  | 分页参数                                                      |

Note: 参数`bk_host_list`和`bk_host_id`只能有一个生效，`bk_host_id`不建议再使用。

#### selectors

| 参数名称     | 参数类型   | 必选 | 描述                                              |
|----------|--------|----|-------------------------------------------------|
| key      | string | 否  | 字段名                                             |
| operator | string | 否  | operator可选值: `=`,`!=`,`exists`,`!`,`in`,`notin` |
| values   | -      | 否  | 不同的operator对应不同的value格式                         |

#### page 字段说明

| 参数名称  | 参数类型 | 必选 | 描述            |
|-------|------|----|---------------|
| start | int  | 是  | 记录开始位置        |
| limit | int  | 是  | 每页限制条数,最大1000 |

### 调用示例

```json
{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 10,
  },
  "bk_module_id": 8,
  "bk_host_list": [11,12],
  "service_instance_ids": [49],
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  }]
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
        "bk_biz_id": 1,
        "id": 49,
        "name": "p1_81",
        "service_template_id": 50,
        "bk_host_id": 11,
        "bk_module_id": 56,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-07-22T09:54:50.906+08:00",
        "last_time": "2019-07-22T09:54:50.906+08:00",
        "bk_supplier_account": "0",
        "service_category_id": 22,
        "process_instances": [
          {
            "process": {
              "proc_num": 0,
              "stop_cmd": "",
              "restart_cmd": "",
              "face_stop_cmd": "",
              "bk_process_id": 43,
              "bk_func_name": "p1",
              "work_path": "",
              "priority": 0,
              "reload_cmd": "",
              "bk_process_name": "p1",
              "pid_file": "",
              "auto_start": false,
              "last_time": "2019-07-22T09:54:50.927+08:00",
              "create_time": "2019-07-22T09:54:50.927+08:00",
              "bk_biz_id": 3,
              "start_cmd": "",
              "user": "",
              "timeout": 0,
              "description": "",
              "bk_supplier_account": "0",
              "bk_start_param_regex": "",
              "bind_info": [
                {
                    "enable": true,
                    "ip": "127.0.0.1",
                    "port": "80",
                    "protocol": "1",
                    "template_row_id": 1234
                }
              ]
            },
            "relation": {
              "bk_biz_id": 1,
              "bk_process_id": 43,
              "service_instance_id": 49,
              "process_template_id": 48,
              "bk_host_id": 11,
              "bk_supplier_account": "0"
            }
          }
        ]
      }
    ]
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

| 参数名称  | 参数类型  | 描述   |
|-------|-------|------|
| count | int   | 总数   |
| info  | array | 返回结果 |

#### data.info 字段说明

| 参数名称                       | 参数类型    | 描述                      |
|----------------------------|---------|-------------------------|
| id                         | integer | 服务实例ID                  |
| name                       | array   | 服务实例名称                  |
| service_template_id        | int     | 服务模板ID                  |
| bk_host_id                 | int     | 主机ID                    |
| bk_host_innerip            | string  | 主机IP                    |
| bk_module_id               | integer | 模块ID                    |
| creator                    | string  | 创建人                     |
| modifier                   | string  | 修改人                     |
| create_time                | string  | 创建时间                    |
| last_time                  | string  | 修复时间                    |
| bk_supplier_account        | string  | 供应商ID                   |
| service_category_id        | integer | 服务分类ID                  |
| process_instances          | array   | 进程实例信息                  |
| bk_biz_id                  | int     | 业务ID                    |
| process_instances.process  | object  | 进程实例详情,进程属性字段           |
| process_instances.relation | object  | 进程实例的关联信息,比如主机ID，进程模板ID |

#### data.info.process_instances[x].process 字段说明

| 参数名称                 | 参数类型   | 描述      |
|----------------------|--------|---------|
| auto_start           | bool   | 是否自动拉起  |
| auto_time_gap        | int    | 拉起间隔    |
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

#### data.info.process_instances[x].process.bind_info[n] 字段说明

| 参数名称            | 参数类型   | 描述                |
|-----------------|--------|-------------------|
| enable          | bool   | 端口是否启用            |
| ip              | string | 绑定的ip             |
| port            | string | 绑定的端口             |
| protocol        | string | 使用的协议             |
| template_row_id | int    | 实例化使用的模板行索引，进程内唯一 |

#### data.info.process_instances[x].relation 字段说明

| 参数名称                | 参数类型   | 描述     |
|---------------------|--------|--------|
| bk_biz_id           | int    | 业务id   |
| bk_process_id       | int    | 进程id   |
| service_instance_id | int    | 服务实例id |
| process_template_id | int    | 进程模版id |
| bk_host_id          | int    | 主机id   |
| bk_supplier_account | string | 开发商账号  |
