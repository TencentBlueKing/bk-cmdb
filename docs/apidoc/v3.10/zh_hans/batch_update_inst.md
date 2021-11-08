### 功能描述

批量更新对象实例

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                |  类型       | 必选   |  描述                            |
|---------------------|-------------|--------|----------------------------------|
| bk_obj_id           | string      | 是     | 模型ID                           |
| update              | object array| 是     | 实例被更新的字段及值             |

#### update
| 字段         | 类型   | 必选  | 描述                           |
|--------------|--------|-------|--------------------------------|
| bk_inst_name | string | 否    | 实例名，也可以为其它自定义字段 |
| datas        | object | 是    | 实例被更新的字段取值           |
| inst_id      | int    | 是    | 指明datas 用于更新的具体实例   |

#### datas

**datas 是map类型的对象，key 是实例对应的模型定义的字段，value是字段的取值**


### 请求参数示例

```python
{
    "bk_supplier_account":"0",
    "bk_obj_id":"test",
    "update":[
        {
          "datas":{
            "bk_inst_name":"batch_update"
          },
          "inst_id":46
         }
        ]
}
```


### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}
```
