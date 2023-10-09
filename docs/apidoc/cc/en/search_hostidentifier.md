### Functional description

Query host identity based on criteria

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field| Type| Required| Description       |
| ---- | ---- | ---- | ---------- |
| ip   |  object |no   | Host ip query criteria|
| page | object |no   | Paging query criteria   |

#### ip

| Field        | Type    | Required| Description       |
| ----------- | ------- | ---- | ---------- |
| data        |  array |no   | Host ip list|
| bk_cloud_id | int     | no   | Cloud area ID |

#### page

| Field| Type   | Required| Description                     |
| ----- | ------ | ---- | ------------------------ |
| start | int    | yes | Record start position             |
| limit | int    | yes | Limit the number of bars per page, with a maximum of 500|
| sort  | string |no   | Sort field                 |



### Request Parameters Example

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

### Return Result Example

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

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data

| Field| Type| Description         |
| ----- | ----- | ------------ |
| count | int   | Number of records     |
| info  | array |Host identity data|

#### data.info[n]
| Field                | Type   | Description                                |
| ------------------- | ------ | ----------------------------------- |
| bk_host_id          |  int    | Host ID                              |
| bk_supplier_account | string |Developer account number                          |
| bk_cloud_id         |  int    | Cloud area ID                        |
| bk_host_innerip     |  string |Intranet IP                              |
| bk_os_type          |  string |Operating system type                        |
| associations        |  dict   | Host mainline Association, key is the module ID to which the host belongs|
| process             |  array  |Host proces information                        |


#### data.info[n].associations
| Field              | Type   | Description                                  |
| ----------------- | ------ | ------------------------------------- |
| bk_biz_id         |  int    | The business ID to which the host belongs                      |
| bk_set_id         |  int    | The set ID to which the host belongs                   |
| bk_module_id      |  int    | Module ID to which the host belongs                      |
| layer             |  dict   | Custom hierarchy info                        |

#### data.info[n].associations.layer
| Field         | Type   | Description               |
| ------------ | ------ | ------------------ |
| bk_inst_id   |  int    | Custom hierarchy instance ID   |
| bk_inst_name | string |Custom hierarchy instance name|
| bk_obj_id    |  int    | Custom hierarchy model ID   |
| child        |  dict   | Custom hierarchy info     |


#### data.info[n].process
| Field                 | Type   | Description                                                         |
| -------------------- | ------ | ------------------------------------------------------------ |
| bk_process_id        |  int    | Process ID                                                       |
| bk_process_name      |  string |Process name                                                       |
| bind_ip              |  string |Binding IP:1/2/3/4(1: 127.0.0.1,2: 0.0.0.0, 3: IP of the first intranet,4: First extranet IP)|
| port                 |  string |Host port                                                     |
| protocol             |  enum   | Protocol:1/2(1:tcp, 2: udp)                                       |
| bk_func_id           |  int    | Function ID                                                       |
| bk_func_name         |  string |Process alias                                                     |
| bk_start_param_regex | string |Process start parameters                                                 |
| bind_modules         |  array  |Module array for process binding                                           |
| bind_info            |  array  |Process binding information                                           |



#### data.info [n] .process.bind .process.bind info [x] description

| Field                 | Type   | Description                                                         |
| -------------------- | ------ | ------------------------------------------------------------ |
|enable| bool| Is the port enabled||
|ip| string| Bound ip||
|port| string| Bound port||
|protocol| string| Protocol used||
|template_row_id| int| Template row index used for instantiation, unique in process|