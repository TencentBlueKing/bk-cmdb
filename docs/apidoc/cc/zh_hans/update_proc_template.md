### 功能描述

更新进程模板信息

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| process_template_id            | int  | 否   | 进程模板ID |
| process_property         | object  | 是   | 需要更新的进程模板字段信息 |

#### process_property 可以出现的字段

注解：

as_default_value: 进程的值是否以模板为准
value: 进程的值，不同的字段类型不一样

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
|proc_num| object|  否|  {"value": null, "as_default_value": false}, value 类型是数字|
|stop_cmd|object| 否| {"value": "","as_default_value": false}, value 类型是字符串|
|restart_cmd|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|face_stop_cmd|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|bk_func_name|object|否|{"value": "a7","as_default_value": true}}, value 类型是字符串|
|work_path|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|priority|object|否|{"value": null,"as_default_value": false}, value 类型是数字|
|reload_cmd|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|bk_process_name|object|否|{"value": "a7","as_default_value": true}}, value 类型是字符串|
|pid_file|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|auto_start|object|否|{"value": null,"as_default_value": null}}, value 类型是boolean|
|auto_time_gap|object|否|{"value": null,"as_default_value": false}, value 类型是数字|
|start_cmd|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|bk_func_id|object|否|{"value": "","as_default_value": false} value 类型是字符串|
|user|object|否|{"value": "","as_default_value": false}, value 类型是字符串|
|timeout|object|否|{"value": null,"as_default_value": false},  value 类型是数字|
|description|object|否|{"value": "1","as_default_value": true}}, value 类型是字符串|
|bk_start_param_regex|object|否|{"value": "","as_default_value": false}, value 类型是字符串||
|bind_info|object|否| {"value":[],,"as_default_value": true }, value详细见process_property.bind_info.value[n]|


#### process_property.bind_info.value[n] 可以出现的字段

注意：

修改bind_info的时候，必须先获取原有进程的bind_info的内容， 然后在进程已有的bind_info上修改，将修改后的内容传递给修改结构。

注解：

as_default_value: 进程的值是否以模板为准
value: 进程的值，不同的字段类型不一样

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
|enable|object|否| {"value": false,"as_default_value": true}, value 类型是boolean|
|ip|object|否| {"value": "1","as_default_value": true}, value 类型是字符串||
|port|object|否| {"value": "100","as_default_value": true}, value 类型是字符串||
|protocol|object|否| {"value": "1","as_default_value": true},, value 类型是字符串||
|row_id|int|否| 唯一表示id,新加的行可以设置为空，更新必须保持原值|







### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "process_template_id": 50,
  "process_property": {
    "proc_num": {
        "value": 300,
        "as_default_value": false
      },
      "stop_cmd": {
        "value": "",
        "as_default_value": false
      },
      "restart_cmd": {
        "value": "",
        "as_default_value": false
      },
      "face_stop_cmd": {
        "value": "",
        "as_default_value": false
      },
      "bk_func_name": {
        "value": "p1",
        "as_default_value": true
      },
      "work_path": {
        "value": "",
        "as_default_value": false
      },
      "priority": {
        "value": null,
        "as_default_value": false
      },
      "reload_cmd": {
        "value": "",
        "as_default_value": false
      },
      "bk_process_name": {
        "value": "p1",
        "as_default_value": true
      },
      "pid_file": {
        "value": "",
        "as_default_value": false
      },
      "auto_start": {
        "value": false,
        "as_default_value": false
      },
      "auto_time_gap": {
        "value": null,
        "as_default_value": false
      },
      "start_cmd": {
        "value": "",
        "as_default_value": false
      },
      "bk_func_id": {
        "value": null,
        "as_default_value": false
      },
      "user": {
        "value": "root100",
        "as_default_value": false
      },
      "timeout": {
        "value": null,
        "as_default_value": false
      },
      "description": {
        "value": "",
        "as_default_value": false
      },
      "bk_start_param_regex": {
        "value": "",
        "as_default_value": false
      },
      "bind_info": {
        "value": [
            {
                "enable": {
                    "value": false,
                    "as_default_value": true
                },
                "ip": {
                    "value": "1",
                    "as_default_value": true
                },
                "port": {
                    "value": "100",
                    "as_default_value": true
                },
                "protocol": {
                    "value": "1",
                    "as_default_value": true
                },
                "row_id": 1
            }
        ],
        "as_default_value": true
      }
  }
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "id": 50,
    "bk_process_name": "p1",
    "bk_biz_id": 1,
    "service_template_id": 51,
    "property": {
      "proc_num": {
        "value": 300,
        "as_default_value": false
      },
      "stop_cmd": {
        "value": "",
        "as_default_value": false
      },
      "restart_cmd": {
        "value": "",
        "as_default_value": false
      },
      "face_stop_cmd": {
        "value": "",
        "as_default_value": false
      },
      "bk_func_name": {
        "value": "p1",
        "as_default_value": true
      },
      "work_path": {
        "value": "",
        "as_default_value": false
      },
      "priority": {
        "value": null,
        "as_default_value": false
      },
      "reload_cmd": {
        "value": "",
        "as_default_value": false
      },
      "bk_process_name": {
        "value": "p1",
        "as_default_value": true
      },
      "pid_file": {
        "value": "",
        "as_default_value": false
      },
      "auto_start": {
        "value": false,
        "as_default_value": false
      },
      "auto_time_gap": {
        "value": null,
        "as_default_value": false
      },
      "start_cmd": {
        "value": "",
        "as_default_value": false
      },
      "bk_func_id": {
        "value": null,
        "as_default_value": false
      },
      "user": {
        "value": "root100",
        "as_default_value": false
      },
      "timeout": {
        "value": null,
        "as_default_value": false
      },
      "description": {
        "value": "",
        "as_default_value": false
      },
      "bk_start_param_regex": {
        "value": "",
        "as_default_value": false
      },
      "bind_info": {
        "value": [
            {
                "enable": {
                    "value": false,
                    "as_default_value": true
                },
                "ip": {
                    "value": "1",
                    "as_default_value": true
                },
                "port": {
                    "value": "100",
                    "as_default_value": true
                },
                "protocol": {
                    "value": "1",
                    "as_default_value": true
                },
                "row_id": 1
            }
        ],
        "as_default_value": true
      }
    },
    "creator": "admin",
    "modifier": "admin",
    "create_time": "2019-06-19T15:24:04.763+08:00",
    "last_time": "2019-06-21T16:25:03.962512+08:00",
    "bk_supplier_account": "0"
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
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data | object | 更新后的进程模板信息 |

#### data 字段说明

| 名称  | 类型  | 描述 |
|---|---|---|
| id | int | 数据id |
| bk_process_name | string |进程别名  |
| bk_biz_id |  int| 业务id |
| service_template_id | int |服务模版id  |
| property |object  |  属性|
| creator              | string             | 本条数据创建者                                                                                 |
| modifier             | string             | 本条数据的最后修改人员            |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| bk_supplier_account | string       | 开发商账号 |
