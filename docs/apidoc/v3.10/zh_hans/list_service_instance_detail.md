### 功能描述

根据业务id查询服务实例列表(带进程信息)，可再附加上模块id等查询条件

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            | int  | 是   | 业务ID |
| bk_module_id         | int  | 否   | 模块ID |
| bk_host_id           | int  | 否   | 主机ID |
| service_instance_ids | int  | 否   | 服务实例ID列表 |
| selectors            | int  | 否   | label过滤功能，operator可选值: `=`,`!=`,`exists`,`!`,`in`,`notin`|
| page                 | object  | 是   | 分页参数 |

#### page 字段说明

| 字段  | 类型   | 必选 | 描述                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | 是   | 记录开始位置          |
| limit | int    | 是   | 每页限制条数,最大1000 |

### 请求参数示例

```python

{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 10,
  },
  "bk_module_id": 8,
  "bk_host_id": 11,
  "service_instance_ids": [49],
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  }]
}


```

### 返回结果示例

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

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | object | 请求返回的数据 |

#### data 字段说明

| 字段|类型|说明|描述|
|---|---|---|---|
|count|integer|总数||
|info|array|返回结果||

#### info 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|integer|服务实例ID||
|name|array|服务实例名称||
|service_template_id|integer|服务模板ID||
|bk_host_id|integer|主机ID||
|bk_host_innerip|string|主机IP||
|bk_module_id|integer|模块ID||
|creator|string|创建人||
|modifier|string|修改人||
|create_time|string|创建时间||
|last_time|string|修复时间||
|bk_supplier_account|string|供应商ID||
|service_category_id|integer|服务分类ID||
|process_instances|数组|进程实例信息|包括||
|bk_biz_id|int|业务ID|业务ID||
|process_instances.process|object|进程实例详情|进程属性字段||
|process_instances.relation|object|进程实例的关联信息|比如主机ID，进程模板ID||

