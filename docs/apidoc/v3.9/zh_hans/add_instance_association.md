### 功能描述

新增模型实例之间的关联关系.

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 是否必填	   |  描述          |
|----------------------|------------|--------|-----------------------------|
| metadata           | object     | Yes    | meta data             |
| condition | string map     | Yes   | 查询条件 |


metadata params

| 字段                 |  类型      | 是否必填	   |  描述         |
|---------------------|------------|--------|-----------------------------|
| label           | string map     | Yes     |标签信息 |


label params

| 字段                 |  类型      | 是否必填	   |  描述         |
|---------------------|------------|--------|-----------------------------|
| bk_biz_id           | string      | Yes     | 业务id |


condition params

| 字段                 |  类型      | 是否必填	   |  描述         |
|---------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           | string     | Yes     | 模型之间关系关系的唯一id|
| bk_inst_id           | int64     | Yes     | 源模型实例id|
| bk_asst_inst_id           | int64     | Yes     | 目标模型实例id|


### 请求参数示例

``` json
{
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 11,
    "bk_asst_inst_id": 21,
    "metadata":{
        "label":{
            "bk_biz_id":"1"
        }
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": 1038
    }
}

```

### 返回结果参数说明

#### data

| 字段       | 类型     | 描述         |
|------------|----------|--------------|
|id|int64|新增的实例关联关系身份id|

