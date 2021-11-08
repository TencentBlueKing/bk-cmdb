### 功能描述

查询实例关联拓扑

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                | 类型   | 必选 | 描述 |
| ------------------- | ------ | ---- | ---- |
| bk_obj_id           | string | 是   | 无   |
| bk_inst_id          | int    | 是   | 无   |


### 请求参数示例

``` python
{
    "bk_supplier_account":"0",
    "bk_obj_id":"test",
    "bk_inst_id":1
}
```


### 返回结果示例

```python
{
    "result":true,
    "code":0,
    "message":"",
    "data":[
        {
            "bk_inst_id":0,
            "bk_inst_name":"",
            "bk_obj_icon":"icon-cc-business",
            "bk_obj_id":"biz",
            "bk_obj_name":"业务",
            "count":1,
            "children":[
                {
                    "bk_inst_id":2,
                    "bk_inst_name":"蓝鲸",
                    "bk_obj_icon":"",
                    "bk_obj_id":"biz",
                    "bk_obj_name":"业务"
                }
            ]
        }
    ]
}
```

### 返回结果参数说明

#### data

| 字段         | 类型         | 描述                           |
| ------------ | ------------ | ------------------------------ |
| bk_inst_id   | int          | 实例ID                         |
| bk_inst_name | string       | 实例用于展示的名字             |
| bk_obj_icon  | string       | 模型图标的名字                 |
| bk_obj_id    | string       | 模型ID                         |
| bk_obj_name  | string       | 模型用于展示的名字             |
| children     | object array | 本模型下所有被关联的实例的集合 |
| count        | int          | children 包含节点的数量        |

#### children

| 字段        | 类型   | 描述               |
|-------------|--------|--------------------|
|bk_inst_id   | int    | 实例ID            |
|bk_inst_name | string | 实例用于展示的名字 |
|bk_obj_icon  | string | 模型图标的名字     |
|bk_obj_id    | string | 模型ID             |
|bk_obj_name  | string | 模型用于展示的名字 |
