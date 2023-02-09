### 功能描述

删除业务集(v3.10.12+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_set_ids      | array     | 是     | 业务集ID列表 |

### 请求参数示例

```python
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_set_ids":[
        10,
        12
    ]
}
```

### 返回结果示例

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission":null,
    "data": {},
    "request_id": "dsda1122adasadadada2222"
}
```
### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| data | object | 请求返回的数据 |
| request_id    | string | 请求链id    |
