### 功能描述

根据服务实例id和设置的标签为服务实例添加标签

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
|instance_ids|array|是|无|服务实例ID|
|labels|object|是|无|添加的Label|

#### labels 字段说明
- key 校验规则: `^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- value 校验规则: `^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### 请求参数示例

```python
{
  "bk_biz_id": 1,
  "instance_ids": [59, 62],
  "labels": {
    "key1": "value1",
    "key2": "value2"
  }
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

### 返回结果说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
