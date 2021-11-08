### 功能描述

根据业务id和服务实例id，以及需要移除的标签，将指定业务下的服务实例移除标签

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            | int  | 是   | 业务ID |
| instance_ids            | array  | 是   | 服务实例ID列表 |
| keys            | array  | 是   | 需要移除的标签key列表 |


### 请求参数示例

```python
{
  "bk_biz_id": 1,
  "instance_ids": [60, 62],
  "keys": ["value1", "value3"]
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
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
