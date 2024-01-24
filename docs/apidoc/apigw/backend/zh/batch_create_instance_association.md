### 描述

批量创建通用模型实例关联关系(版本：v3.10.2+，权限：源模型实例和目标模型实例的编辑权限)

### 输入参数

| 参数名称           | 参数类型   | 必选 | 描述                     |
|----------------|--------|----|------------------------|
| bk_obj_id      | string | 是  | 源模型id                  |
| bk_asst_obj_id | string | 是  | 目标模型模型id               |
| bk_obj_asst_id | string | 是  | 模型之间关系关系的唯一id          |
| details        | array  | 是  | 批量创建关联关系的内容，不能超过200个关系 |

#### details

| 参数名称            | 参数类型 | 必选 | 描述       |
|-----------------|------|----|----------|
| bk_inst_id      | int  | 是  | 源模型实例id  |
| bk_asst_inst_id | int  | 是  | 目标模型实例id |

### 调用示例

```json
{
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

### 响应示例

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission": null,
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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称            | 参数类型 | 描述                                                |
|-----------------|------|---------------------------------------------------|
| success_created | map  | key为实例关联关系在参数details数组中的index，value为创建成功的实例关联关系id |
| error_msg       | map  | key为实例关联关系在参数details数组中的index，value为失败信息          |
