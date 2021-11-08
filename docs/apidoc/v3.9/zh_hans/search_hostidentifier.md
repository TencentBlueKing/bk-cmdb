### 功能描述

根据条件查询主机身份

### 请求参数


#### 通用参数

| 字段 | 类型 | 必选 |  描述 |
|-----------|------------|--------|------------|
| bk_app_code  |  string    | 是 | 应用ID     |
| bk_app_secret|  string    | 是 | 安全密钥(应用 TOKEN)，可以通过 蓝鲸智云开发者中心 -&gt; 点击应用ID -&gt; 基本信息 获取 |
| bk_token     |  string    | 否 | 当前用户登录态，bk_token与bk_username必须一个有效，bk_token可以通过Cookie获取 |
| bk_username  |  string    | 否 | 当前用户用户名，应用免登录态验证白名单中的应用，用此字段指定当前用户 |

#### 接口参数

| 字段 | 类型 | 必选 | 描述       |
| ---- | ---- | ---- | ---------- |
| ip   | dict | 否   | 主机ip查询条件 |
| page | dict | 否   | 分页查询条件   |

#### ip

| 字段        | 类型    | 必选 | 描述       |
| ----------- | ------- | ---- | ---------- |
| data        | array | 否   | 主机ip列表 |
| bk_cloud_id | int     | 否   | 云区域ID   |

#### page

| 字段  | 类型   | 必选 | 描述                     |
| ----- | ------ | ---- | ------------------------ |
| start | int    | 是   | 记录开始位置             |
| limit | int    | 是   | 每页限制条数,最大值为500 |
| sort  | string | 否   | 排序字段                 |



### 请求参数示例

```json
{
    "ip": {
        "data": [
            "192.168.0.1"
        ],
        "bk_cloud_id": 0
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"bk_host_name"
    }
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
        "count": 1,
        "info": [
            {
                "bk_host_id": 4,
                "bk_host_name": "",
                "bk_supplier_account": "",
                "bk_cloud_id": 0,
                "bk_cloud_name": "default area",
                "bk_host_innerip": "192.168.0.1",
                "bk_host_outerip": "",
                "bk_os_type": "",
                "bk_os_name": "",
                "bk_mem": 0,
                "bk_cpu": 0,
                "bk_disk": 0,
                "associations": {
                    "51": {
                        "bk_biz_id": 3,
                        "bk_biz_name": "test",
                        "bk_set_id": 8,
                        "bk_set_name": "test",
                        "bk_module_id": 51,
                        "bk_module_name": "test",
                        "bk_service_status": "1",
                        "bk_set_env": "3",
                         "layer": {
                              "bk_inst_id": 3,
                              "bk_inst_name": "a",
                              "bk_obj_id": "a",
                              "child": {
                                    "bk_inst_id": 5,
                                    "bk_inst_name": "b",
                                    "bk_obj_id": "b",
                                    "child": {}
                              }
                        }
                    }
                },
                "process": [
                    {
                        "bk_process_id": 43,
                        "bk_process_name": "test",
                        "bind_ip": "",
                        "port": "8000",
                        "protocol": "1",
                        "bk_func_id": "",
                        "bk_func_name": "test",
                        "bk_start_param_regex": "",
                        "bind_info": [
                            {
                                "enable": false,  
                                "ip": "127.0.0.1",  
                                "port": "8000",  
                                "protocol": "1", 
                                "template_row_id": 1  
                            }
                        ],
                        "bind_modules": [
                            51
                        ]
                    }
                ]
            }
        ]
    }
}
```

### 返回结果参数说明

#### data

| 字段  | 类型  | 描述         |
| ----- | ----- | ------------ |
| count | int   | 记录条数     |
| info  | array | 主机身份数据 |

#### data.info[n]
| 字段                | 类型   | 描述                                |
| ------------------- | ------ | ----------------------------------- |
| bk_host_id          | int    | 主机ID                              |
| bk_host_name        | string | 主机名称                            |
| bk_supplier_account | string | 开发商账号                          |
| bk_cloud_id         | int    | 云区域ID                            |
| bk_cloud_name       | string | 云区域名称                          |
| bk_host_innerip     | string | 内网IP                              |
| bk_host_outerip     | string | 外网IP                              |
| bk_os_type          | string | 操作系统类型                        |
| bk_os_name          | string | 操作系统名称                        |
| bk_mem              | string | 内存容量                            |
| bk_cpu              | int    | CPU逻辑核心数                       |
| bk_disk             | int    | 磁盘容量                            |
| associations        | dict   | 主机主线关联，key为主机所属的模块ID |
| process             | array  | 主机进程信息                        |


#### data.info[n].associations
| 字段              | 类型   | 描述                                  |
| ----------------- | ------ | ------------------------------------- |
| bk_biz_id         | int    | 主机所属的业务ID                      |
| bk_biz_name       | string | 主机所属的业务名字                    |
| bk_set_id         | int    | 主机所属的集群ID                      |
| bk_set_name       | string | 主机所属的集群名字                    |
| bk_module_id      | int    | 主机所属的模块ID                      |
| bk_module_name    | string | 主机所属的模块名字                    |
| bk_service_status | enum   | 服务状态:1/2(1:开放,2:关闭)           |
| bk_set_env        | enum   | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
| layer             | dict   | 自定义层级信息                        |

#### data.info[n].associations.layer
| 字段         | 类型   | 描述               |
| ------------ | ------ | ------------------ |
| bk_inst_id   | int    | 自定义层级实例ID   |
| bk_inst_name | string | 自定义层级实例名字 |
| bk_obj_id    | int    | 自定义层级模型ID   |
| child        | dict   | 自定义层级信息     |


#### data.info[n].process
| 字段                 | 类型   | 描述                                                         |
| -------------------- | ------ | ------------------------------------------------------------ |
| bk_process_id        | int    | 进程ID                                                       |
| bk_process_name      | string | 进程名                                                       |
| bind_ip              | string | 绑定IP:1/2/3/4(1:127.0.0.1,2:0.0.0.0,3:第一内网IP,4:第一外网IP) |
| port                 | string | 主机端口                                                     |
| protocol             | enum   | 协议:1/2(1:tcp, 2:udp)                                       |
| bk_func_id           | int    | 功能ID                                                       |
| bk_func_name         | string | 进程别名                                                     |
| bk_start_param_regex | string | 进程启动参数                                                 |
| bind_modules         | array  | 进程绑定的模块数组                                           |
| bind_info            | array  | 进程绑定信息                                           |



#### data.info[n].process.bind_info[x] 描述

| 字段                 | 类型   | 描述                                                         |
| -------------------- | ------ | ------------------------------------------------------------ |
|enable|bool|端口是否启用||
|ip|string|绑定的ip||
|port|string|绑定的端口||
|protocol|string|使用的协议||
|row_id|int|实例化使用的模板行索引，进程内唯一|