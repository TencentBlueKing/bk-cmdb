### 功能描述

根据条件查询主机身份

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段 | 类型 | 必选 | 描述       |
| ---- | ---- | ---- | ---------- |
| ip   | object | 否   | 主机ip查询条件 |
| page | object | 否   | 分页查询条件   |

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
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "data": {
        "count": 1,
        "info": [
            {
                "bk_host_id": 11,
                "bk_cloud_id": 0,
                "bk_host_innerip": "1.1.1.1",
                "bk_os_type": "",
                "bk_supplier_account": "0",
                "associations": {
                    "15553": {
                        "bk_biz_id": 11,
                        "bk_set_id": 4760,
                        "bk_module_id": 15553,
                        "layer": null
                    }
                },
                "process": [
                    {
                        "bk_process_id": 90908,
                        "bk_process_name": "test",
                        "bind_ip": "1.1.1.1",
                        "port": "8080",
                        "protocol": "1",
                        "bk_func_name": "test",
                        "bk_start_param_regex": "./test",
                        "bk_enable_port": true,
                        "bind_info": [
                            {
                                "enable": false,
                                "ip": "1.1.1.1",
                                "port": "8080",
                                "protocol": "1",
                                "template_row_id": 1
                            },
                            {
                                "enable": true,
                                "ip": "127.0.0.1",
                                "port": "8081",
                                "protocol": "2",
                                "template_row_id": 2
                            }
                        ]
                    }
                ]
            }
        ]
    },
    "message": "success",
    "permission": null,
    "request_id": "cc26632ed3c344c79c0002ae9bcf3009"
}
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

#### data

| 字段  | 类型  | 描述         |
| ----- | ----- | ------------ |
| count | int   | 记录条数     |
| info  | array | 主机身份数据 |

#### data.info[n]
| 字段                | 类型   | 描述                                |
| ------------------- | ------ | ----------------------------------- |
| bk_host_id          | int    | 主机ID                              |
| bk_supplier_account | string | 开发商账号                          |
| bk_cloud_id         | int    | 云区域ID                            |
| bk_host_innerip     | string | 内网IP                              |
| bk_os_type          | string | 操作系统类型                        |
| associations        | dict   | 主机主线关联，key为主机所属的模块ID |
| process             | array  | 主机进程信息                        |


#### data.info[n].associations
| 字段              | 类型   | 描述                                  |
| ----------------- | ------ | ------------------------------------- |
| bk_biz_id         | int    | 主机所属的业务ID                      |
| bk_set_id         | int    | 主机所属的集群ID                      |
| bk_module_id      | int    | 主机所属的模块ID                      |
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
|template_row_id|int|实例化使用的模板行索引，进程内唯一|