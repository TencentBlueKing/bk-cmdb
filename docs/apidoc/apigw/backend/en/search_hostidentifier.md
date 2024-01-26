### Description

Query host identity based on conditions

### Parameters

| Name | Type   | Required | Description                 |
|------|--------|----------|-----------------------------|
| ip   | object | No       | Host IP query conditions    |
| page | object | No       | Pagination query conditions |

#### ip

| Name        | Type  | Required | Description     |
|-------------|-------|----------|-----------------|
| data        | array | No       | Host IP list    |
| bk_cloud_id | int   | No       | Control area ID |

#### page

| Name  | Type   | Required | Description                           |
|-------|--------|----------|---------------------------------------|
| start | int    | Yes      | Record start position                 |
| limit | int    | Yes      | Each page limit, maximum value is 500 |
| sort  | string | No       | Sorting field                         |

### Request Example

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

### Response Example

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
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name  | Type  | Description        |
|-------|-------|--------------------|
| count | int   | Record count       |
| info  | array | Host identity data |

#### data.info[n]

| Name                | Type   | Description                                                               |
|---------------------|--------|---------------------------------------------------------------------------|
| bk_host_id          | int    | Host ID                                                                   |
| bk_supplier_account | string | Developer account                                                         |
| bk_cloud_id         | int    | Control area ID                                                           |
| bk_host_innerip     | string | Intranet IP                                                               |
| bk_os_type          | string | Operating system type                                                     |
| associations        | dict   | Host mainline association, key is the module ID to which the host belongs |
| process             | array  | Host process information                                                  |

#### data.info[n].associations

| Name         | Type | Description                          |
|--------------|------|--------------------------------------|
| bk_biz_id    | int  | Business ID of the host              |
| bk_set_id    | int  | Cluster ID to which the host belongs |
| bk_module_id | int  | Module ID to which the host belongs  |
| layer        | dict | Custom level information             |

#### data.info[n].associations.layer

| Name         | Type   | Description                |
|--------------|--------|----------------------------|
| bk_inst_id   | int    | Custom level instance ID   |
| bk_inst_name | string | Custom level instance name |
| bk_obj_id    | int    | Custom level model ID      |
| child        | dict   | Custom level information   |

#### data.info[n].process

| Name                 | Type   | Description                                                                     |
|----------------------|--------|---------------------------------------------------------------------------------|
| bk_process_id        | int    | Process ID                                                                      |
| bk_process_name      | string | Process name                                                                    |
| bind_ip              | string | Bound IP:1/2/3/4(1:127.0.0.1,2:0.0.0.0,3:First Intranet IP,4:First External IP) |
| port                 | string | Host port                                                                       |
| protocol             | enum   | Protocol:1/2(1:tcp, 2:udp)                                                      |
| bk_func_id           | int    | Function ID                                                                     |
| bk_func_name         | string | Process alias                                                                   |
| bk_start_param_regex | string | Process startup parameters                                                      |
| bind_modules         | array  | Array of modules bound by the process                                           |
| bind_info            | array  | Process binding information                                                     |

#### data.info[n].process.bind_info[x] Description

| Name            | Type   | Description                                                          |
|-----------------|--------|----------------------------------------------------------------------|
| enable          | bool   | Whether the port is enabled                                          |
| ip              | string | Bound IP                                                             |
| port            | string | Bound port                                                           |
| protocol        | string | Used protocol                                                        |
| template_row_id | int    | Template row index used for instantiation, unique within the process |
