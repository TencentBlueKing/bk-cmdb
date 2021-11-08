### Functional description

list service instances with processes info

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| bk_biz_id            | int  | Yes  | Business ID |
| bk_module_id         | int  | No   | Module ID |
| bk_host_id           | int  | No   | Host ID, deprecated: please do not use any more |
| bk_host_list         | array| No   | Host ID list |
| service_instance_ids | int  | No   | Service Instance IDs |
| selectors            | int  | No   | label filters，available operator values are: `=`,`!=`,`exists`,`!`,`in`,`notin`|
| page                 | object| Yes | page paremeters |

Only one parameter between `bk_host_list` and `bk_host_id` can take effect. `bk_host_id` does not recommend using it again.
#### page params

| Field                 |  Type      | Required	   |  Description       | 
|--------|------------|--------|------------|
|start|int|No|get the data offset location|
|limit|int|Yes|The number of data points in the past is limited, suggest 1000|

### Request Parameters Example

```python

{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 1
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

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |

#### Data field description

| Field       | Type     | Description         |
|---|---|---|---|
|count|integer|total count||
|info|array|response data||

#### Info field description

| Field       | Type     | Description         |
|---|---|---|---|
|id|integer|Service Instance ID||
|name|array|Service Instance Name||
|service_template_id|integer|Service Template ID||
|bk_host_id|integer|Host ID||
|bk_host_innerip|string|Host IP||
|bk_module_id|integer|Module ID||
|creator|string|Creator||
|modifier|string|Modifier||
|create_time|string|Create Time||
|last_time|string|Update Time||
|bk_supplier_account|string|Supplier Account ID||
|process_instances|Array|Process Instance Data|||
|process_instances.process|object|Process Instance Detail|Process Instance Property||
|process_instances.relation|object|Process Instance Relations|f.e. host id, process template id||

