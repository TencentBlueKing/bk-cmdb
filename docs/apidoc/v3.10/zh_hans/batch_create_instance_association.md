### 功能描述

 批量创建通用模型实例关联关系(v3.10.2+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 参数           | 类型   | 必选 | 描述                     |
| -------------- | ------ | ---- | ------------------------ |
| bk_obj_id      | string | 是   | 源模型id                 |
| bk_asst_obj_id | string | 是   | 目标模型模型id           |
| bk_obj_asst_id | string | 是   | 模型之间关系关系的唯一id |
| details        | array  | 是   | 批量创建关联关系的内容，不能超过200个关系        |

#### details

| 参数            | 类型   | 必选 | 描述           |
| --------------- | ------ | ---- | -------------- |
| bk_inst_id      | int | 是   | 源模型实例id   |
| bk_asst_inst_id | int | 是   | 目标模型实例id |

#### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"bk_switch",
    "bk_asst_obj_id":"host",
    "bk_obj_asst_id":"bk_switch_belong_host",
    "details":[
        {
            "bk_inst_id":11,
            "bk_asst_inst_id":21
        },
        {
            "bk_inst_id":12,
            "bk_asst_inst_id":22
        }
    ]
}
```

### 返回结果示例

```json
{
    "result":true,
    "code":0,
    "message":"",
    "data":{
        "success_created":{
            "0":73
        },
        "error_msg":{
            "1":"关联实例不存在"
        }
    }
}
```

### 返回结果参数说明

#### data

| 字段            | 类型 | 描述                                                     |
| -------------- | ---- | -------------------------------------------------------- |
| success_created | map | key为实例关联关系在参数details数组中的index，value为创建成功的实例关联关系id |
| error_msg       | map | key为实例关联关系在参数details数组中的index，value为失败信息          |