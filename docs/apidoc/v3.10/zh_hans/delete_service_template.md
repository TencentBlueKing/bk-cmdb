### 功能描述

根据服务模板ID删除服务模板

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| service_template_id | int  | 是   | 服务模板ID |

### 请求参数示例

```python
{
  "bk_biz_id": 1,
  "service_template_id": 1
}
```

### 返回结果示例

```python
{
  "result": false,
  "code": 1199056,
  "message": "禁止删除被引用记录",
  "data": ""
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
