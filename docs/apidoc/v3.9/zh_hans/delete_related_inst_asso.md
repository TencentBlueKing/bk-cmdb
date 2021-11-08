### 功能描述

 根据实例关联关系的ID删除实例之间的关联。(生效版本：3.5.40)

### 请求参数

{{ common_args_desc }}

 
#### 接口参数

| 字段 | 类型     | 必选 | 描述             |
| :--- | :------- | :--- | :--------------- |
| id   | int | 是   | 实例关联关系的ID（注：非模型实例的身份ID）， 最多500个 |
| bk_obj_id | string | 是 | 关联关系源模型的模型唯一名称 |

### 请求参数示例

```json
{
    "id":[1,2],
    "bk_obj_id": "abc"
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### 返回结果参数说明

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | object | 请求返回的数据 |