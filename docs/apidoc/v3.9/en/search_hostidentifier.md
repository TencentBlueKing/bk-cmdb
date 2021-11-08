### Functional description

search host identifier

### General Parameters

{{ common_args_desc }}

#### Request Parameters

| Field | Type | Required | Description             |
| ----- | ---- | -------- | ----------------------- |
| ip    | dict | No       | host ip query condition |
| page  | dict | No       | paging query condition  |

#### ip

| Field       | Type  | Required | Description   |
| ----------- | ----- | -------- | ------------- |
| data        | array | No       | host ip list  |
| bk_cloud_id | int   | No       | cloud area ID |

#### page

| Field | Type   | Required | Description            |
| ----- | ------ | -------- | ---------------------- |
| start | int    | Yes      | start record           |
| limit | int    | Yes      | page limit, max is 500 |
| sort  | string | No       | the field for sort     |



### Request Parameters Example

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

### Return Result Example

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
                        "port": "8000,8001,8003,8004",
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

### Return Result Parameters Description

#### data

| Field | Type  | Description          |
| ----- | ----- | -------------------- |
| count | int   | the num of record    |
| info  | array | host identifier data |

#### data.info
| Field               | Type   | Description                                                  |
| ------------------- | ------ | ------------------------------------------------------------ |
| bk_host_id          | int    | host ID                                                      |
| bk_host_name        | string | host name                                                    |
| bk_supplier_account | string | supplier account                                             |
| bk_cloud_id         | int    | cloud area ID                                                |
| bk_cloud_name       | string | cloud area name                                              |
| bk_host_innerip     | string | inner ip                                                     |
| bk_host_outerip     | string | outer ip                                                     |
| bk_os_type          | string | os type                                                      |
| bk_os_name          | string | os name                                                      |
| bk_mem              | string | memory capacity                                              |
| bk_cpu              | int    | CPU count                                                    |
| bk_disk             | int    | disk capacity                                                |
| associations        | dict   | host mainline associations, key is the module ID that host belongs to |
| process             | array  | host's process info                                          |


#### data.info.associations
| Field             | Type   | Description                                          |
| ----------------- | ------ | ---------------------------------------------------- |
| bk_biz_id         | int    | host biz ID                                          |
| bk_biz_name       | string | host biz name                                        |
| bk_set_id         | int    | host set ID                                          |
| bk_set_name       | string | host set name                                        |
| bk_module_id      | int    | host module ID                                       |
| bk_module_name    | string | host module name                                     |
| bk_service_status | enum   | the service status:1/2 (1:open,2:close)              |
| bk_set_env        | enum   | environment type:1/2/3(1:test,2:experience,3:formal) |
| layer             | dict   | self-defined layer info                              |

#### data.info.associations.layer
| Field        | Type   | Description                  |
| ------------ | ------ | ---------------------------- |
| bk_inst_id   | int    | self-defined layer inst ID   |
| bk_inst_name | string | self-defined layer inst name |
| bk_obj_id    | int    | self-defined layer modle ID  |
| child        | dict   | self-defined layer info      |


#### data.info.process
| Field                | Type   | Description                                                  |
| -------------------- | ------ | ------------------------------------------------------------ |
| bk_process_id        | int    | process ID                                                   |
| bk_process_name      | string | process name                                                 |
| bind_ip              | object | bind IP:1/2/3/4(1:127.0.0.1,2:0.0.0.0,3:first intranet IP,4:first extranet IP) |
| port                 | string | host port                                                    |
| protocol             | enum   | protocol:1/2(1:tcp, 2:udp)                                   |
| bk_func_id           | int    | process function ID                                          |
| bk_func_name         | string | process alias                                                |
| bk_start_param_regex | string | process start parameter                                      |
| bind_modules         | array  | process module array                                         |
| bind_info            | array  | process bind info                                            |


#### data.info[n].process.bind_info[x] Description
| Field       | Type     | Description         |
|---|---|---|---|
|enable|bool|Whether the port is enabled||
|ip|string|bind ip||
|port|string|bind port||
|protocol|string|protocol used||
|row_id|int|template row index used for instantiation, unique in the process|