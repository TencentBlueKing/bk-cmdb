### 功能描述

查询模型的实例关联关系。

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 是否必填	   |  描述          |
|----------------------|------------|--------|-----------------------------|
| condition | string map     | Yes   | 查询条件 |
| bk_obj_id           | string     | YES     | 源模型id(v3.10+)|


condition params

| 字段                 |  类型      | 是否必填	   |  描述         |
|---------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           | string     | Yes     | 模型关联关系的唯一id|
| bk_asst_id           | string     | NO     | 关联类型的唯一id|
| bk_asst_obj_id           | string     | NO     | 目标模型id|


### 请求参数示例

``` json
{
    "condition": {
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_asst_id": "",
        "bk_asst_obj_id": ""
    },
    "bk_object_id": "xxx"
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": [{
        "bk_obj_asst_id": "bk_switch_belong_bk_host",
        "bk_obj_id":"switch",
        "bk_asst_obj_id":"host",
        "bk_inst_id":12,
        "bk_asst_inst_id":13
    }]
}

```


### 返回结果参数说明

#### data

| 字段       | 类型     | 描述         |
|------------|----------|--------------|
|id|int64|the association's unique id|
| bk_obj_asst_id| string|  自动生成的模型关联关系id.|
| bk_obj_id| string| 关联关系源模型id |
| bk_asst_obj_id| string| 关联关系目标模型id|
| bk_inst_id| int64| 源模型实例id|
| bk_asst_inst_id| int64| 目标模型实例id|

