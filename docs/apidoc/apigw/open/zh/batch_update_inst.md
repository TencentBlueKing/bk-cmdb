### 描述

批量更新对象实例(权限：模型实例编辑权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述         |
|-----------|--------|----|------------|
| bk_obj_id | string | 是  | 模型ID       |
| update    | array  | 是  | 实例被更新的字段及值 |

#### update

| 参数名称    | 参数类型   | 必选 | 描述                |
|---------|--------|----|-------------------|
| datas   | object | 是  | 实例被更新的字段取值        |
| inst_id | int    | 是  | 指明datas 用于更新的具体实例 |

#### datas

| 参数名称         | 参数类型   | 必选 | 描述              |
|--------------|--------|----|-----------------|
| bk_inst_name | string | 否  | 实例名，也可以为其它自定义字段 |

**datas 是map类型的对象，key 是实例对应的模型定义的字段，value是字段的取值**

### 调用示例

```json
{
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

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": "success"
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
