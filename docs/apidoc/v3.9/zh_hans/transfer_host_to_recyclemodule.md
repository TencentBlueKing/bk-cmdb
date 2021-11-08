### 功能描述

上交主机到业务的待回收模块

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | 否     | 开发商账号 |
| bk_biz_id     |  int     | 是     | 业务ID |
| bk_host_id    |  array   | 是     | 主机ID |

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "bk_biz_id": 1,
    "bk_host_id": [
        9,
        10
    ]
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null
}
```
