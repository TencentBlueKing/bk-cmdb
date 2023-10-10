### Functional description

Query the service instance list (with process information) according to the service id, and query conditions such as module id can be added

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            |  int  |yes   | Business ID |
| bk_module_id         |  int  |no   | Module ID|
| bk_host_id           |  int  |no   | Host ID, Note: This field is no longer maintained, please use bk_host_list field|
| bk_host_list         |  array| no   | Host ID list|
| service_instance_ids | int  |no   | Service instance ID list|
| selectors            |  int  |no   | Label filtering function, operator optional value: `=`,`!=`,` exists`,`!`,` in`,`notin`|
| page                 |  object  |yes   | Paging parameter|

Note: only one of the parameters`bk_host_list` and`bk_host_id` can be effective`bk_host_id`. It is not recommended to use it again.
#### page params

| Field                 |  Type      | Required	   |  Description       | 
|--------|------------|--------|------------|
|start|int|No|get the data offset location|
|limit|int|Yes|page limit, maximum value is 1000|
#### selectors
| Field                 | Type      | Required	   | Description                 |
| -------- | ------ | ---- | ------ |
| key    |  string |no   | Field name|
| operator | string |no   | Operator optional value: `=`,`!=`,` exists`,`!`,` in`,`notin` |
| values    | -      |no| Different values correspond to different value formats                            |

#### Page field Description

| Field| Type   | Required| Description                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | yes | Record start position          |
| limit | int    | yes | Limit bars per page, Max. 1000|

### Request Parameters Example

```python

{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
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

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

#### Data field Description

| Field| Type| Description|
|---|---|---|
|count| int| Total|
|info| array| Return result|

#### Data.info Field Description

| Field| Type| Description|
|---|---|---|
|id| integer| Service instance ID||
|name| array| Service instance name||
|service_template_id| int| Service template ID||
|bk_host_id| int| Host ID||
|bk_host_innerip| string| Host IP||
|bk_module_id| integer| Module ID||
|creator| string| Founder||
|modifier| string| Modified by||
|create_time| string| Settling time||
|last_time| string| Repair time||
|bk_supplier_account| string| Vendor ID||
|service_category_id| integer| Service class ID||
|process_instances| Array| Process instance information| Including|
|bk_biz_id| int| Service ID| Business ID |
|process_instances.process| object| Process instance details| Process properties field|
|process_instances.relation| object| Process instance association information| Such as host ID, proces template ID|

#### Data.info.process_instances [x] .process .process description
| Field| Type| Description|
|---|---|---|
|auto_start| bool| Whether to pull up automatically|
|auto_time_gap| int| Pull up interval|
|bk_biz_id| int| Business ID |
|bk_func_id| string| Function ID|
|bk_func_name| string| Process name|
|bk_process_id| int| Process id|
|bk_process_name| string| Process alias|
|bk_start_param_regex| string| Process start parameters|
|bk_supplier_account| string| Developer account number|
|create_time| string| Settling time|
|description| string| Description|
|face_stop_cmd| string| Forced stop command|
|last_time| string| Update time|
|pid_file| string| PID file path|
|priority| int| Startup priority|
|proc_num| int| Number of starts|
|reload_cmd| string| Process reload command|
|restart_cmd| string| Restart command|
|start_cmd| string| Start command|
|stop_cmd| string| Stop order|
|timeout| int| Operation time-out duration|
|user| string| Start user|
|work_path| string| Working path|
|bind_info| object| Binding information|

#### Data.info.process_instances [x] .process.bind .process.bind info [n] Field Description
| Field| Type| Description|
|---|---|---|
|enable| bool| Is the port enabled|
|ip| string| Bound ip|
|port| string| Bound port|
|protocol| string| Protocol used|
|template_row_id| int| Template row index used for instantiation, unique in process|

#### Data.info.process_instances [x]. Relationfield description
| Field| Type| Description|
|---|---|---|
|bk_biz_id| int| Business ID |
|bk_process_id| int| Process id|
|service_instance_id| int| Service instance id|
|process_template_id| int| Process template id|
|bk_host_id| int| Host id|
|bk_supplier_account| string| Developer account number|


