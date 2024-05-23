### 描述

点分五位查询进程实例的相关信息 (v3.9.13)

- 该接口专供GSEKit使用，在ESB文档中为hidden状态

### 输入参数

| 参数名称             | 参数类型         | 必选 | 描述                                                                                        |
|------------------|--------------|----|-------------------------------------------------------------------------------------------|
| bk_biz_id        | int64        | 是  | 业务ID                                                                                      |
| bk_set_ids       | int64 array  | 否  | 集群ID列表，若为空，则代表任意一集群                                                                       |
| bk_module_ids    | int64 array  | 否  | 模块ID列表，若为空，则代表任意一模块                                                                       |
| ids              | int64 array  | 否  | 服务实例ID列表，若为空，则代表任意一实例                                                                     |
| bk_process_names | string array | 否  | 进程名称列表，若为空，则代表任意一进程。                                                                      |
| bk_process_ids   | int64 array  | 否  | 进程ID列表，若为空，则代表任一进程                                                                        |
| fields           | string array | 否  | 进程属性列表，控制返回结果的进程实例信息里有哪些字段，能够加速接口请求和减少网络流量传输<br>为空时返回进程所有字段,bk_process_id,bk_process_name |
| page             | dict         | 是  | 分页条件                                                                                      |

这些字段的条件关系是关系与(&amp;&amp;)，只会查询同时满足所填条件的进程实例<br>
举例来说：如果同时填了bk_set_ids和bk_module_ids，而bk_module_ids都不属于bk_set_ids，则查询结果为空

#### page

| 参数名称  | 参数类型   | 必选 | 描述                                        |
|-------|--------|----|-------------------------------------------|
| start | int    | 否  | 记录开始位置，默认为0                               |
| limit | int    | 是  | 每页限制条数,最大500                              |
| sort  | string | 否  | 排序字段，'-'表示倒序, 只能是进程的字段，默认按bk_process_id排序 |

### 调用示例

```json
{
    "bk_biz_id": 3,
    "set": {
        "bk_set_ids": [
            11,
            12
        ]
    },
    "module": {
        "bk_module_ids": [
            60,
            61
        ]
    },
    "service_instance": {
        "ids": [
            4,
            5
        ]
    },
    "process": {
        "bk_process_names": [
            "pr1",
            "alias_pr2"
        ],
        "bk_process_ids": [
            45,
            46,
            47
        ]
    },
    "fields": [
        "bk_process_id",
        "bk_process_name",
        "bk_func_name"
    ],
    "page": {
        "start": 0,
        "limit": 100,
        "sort": "bk_process_id"
    }
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "count": 2,
        "info": [
            {
                "set": {
                    "bk_set_id": 11,
                    "bk_set_name": "set1",
                    "bk_set_env": "3"
                },
                "module": {
                    "bk_module_id": 60,
                    "bk_module_name": "mm1"
                },
                "host": {
                    "bk_host_id": 4,
                    "bk_cloud_id": 0,
                    "bk_host_innerip": "127.0.0.1",
                    "bk_host_innerip_v6":"1::1",
                    "bk_addressing":"dynamic",
                    "bk_agent_id":"xxxxxx"
                },
                "service_instance": {
                    "id": 4,
                    "name": "127.0.0.1_pr1_3333"
                },
                "process_template": {
                    "id": 48
                },
                "process": {
                    "bk_func_name": "pr1",
                    "bk_process_id": 45,
                    "bk_process_name": "pr1"
                }
            },
            {
                "set": {
                    "bk_set_id": 11,
                    "bk_set_name": "set1",
                    "bk_set_env": "3"
                },
                "module": {
                    "bk_module_id": 60,
                    "bk_module_name": "mm1"
                },
                "host": {
                    "bk_host_id": 4,
                    "bk_cloud_id": 0,
                    "bk_host_innerip": "127.0.0.1"
                },
                "service_instance": {
                    "id": 4,
                    "name": "127.0.0.1_pr1_3333"
                },
                "process_template": {
                    "id": 49
                },
                "process": {
                    "bk_func_name": "pr2",
                    "bk_process_id": 46,
                    "bk_process_name": "alias_pr2"
                }
            }
        ]
    }
}
```

### 响应参数说明

| 参数名称    | 参数类型   | 描述                         |
|---------|--------|----------------------------|
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误  |
| message | string | 请求失败返回的错误信息                |

#### data 字段说明

| 参数名称             | 参数类型   | 描述           |
|------------------|--------|--------------|
| count            | int    | 符合条件的进程实例总数量 |
| set              | object | 进程所属的集群信息    |
| module           | object | 进程所属的模块信息    |
| host             | object | 进程所属的主机信息    |
| service_instance | object | 进程所属的服务实例信息  |
| process_template | object | 进程模板信息       |
| process          | object | 进程自身的详细信息    |

#### data.set 字段说明

| 参数名称        | 参数类型   | 描述   |
|-------------|--------|------|
| bk_set_id   | int    | 集群id |
| bk_set_name | string | 集群名称 |
| bk_set_env  | string | 环境类型 |

#### data.module 字段说明

| 参数名称           | 参数类型   | 描述   |
|----------------|--------|------|
| bk_module_id   | int    | 模块id |
| bk_module_name | string | 模块名称 |

#### data.host 字段说明

| 参数名称               | 参数类型   | 描述       |
|--------------------|--------|----------|
| bk_host_id         | int    | 主机id     |
| bk_cloud_id        | int    | 管控区域id   |
| bk_host_innerip    | string | 主机内网IP   |
| bk_host_innerip_v6 | int    | 主机内网IPv6 |
| bk_addressing      | string | 寻址方式     |
| bk_agent_id        | string | Agent ID |

#### data.service_instance 字段说明

| 参数名称 | 参数类型   | 描述     |
|------|--------|--------|
| id   | int    | 服务实例id |
| name | string | 服务实例名称 |

#### data.process_template 字段说明

| 参数名称 | 参数类型 | 描述     |
|------|------|--------|
| id   | int  | 集群模板id |

#### data.process 字段说明

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
