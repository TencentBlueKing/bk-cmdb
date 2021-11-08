### 功能描述

退订事件

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段               |  类型      | 必选   |  描述      |
|--------------------|------------|--------|------------|
|subscription_id     | int        | 是     | 订阅ID     |

### 请求参数示例

```python
{
    "bk_supplier_account":"0",
    "subscription_id":1
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data":"success"
}
```
