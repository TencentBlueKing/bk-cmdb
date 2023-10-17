### 功能描述

获取主机基础信息详情

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | 否     | 开发商账号 |
| bk_host_id     |  int       | 是     | 主机身份ID，即bk_host_id字段值 |

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_host_id": 10000
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "request_id": "c11aasdadadadsadasdadasd1111ds"
    "permission": null,
    "data": [
        {
            "bk_property_id": "bk_host_innerip",
            "bk_property_name": "内网IP",
            "bk_property_value": "127.0.0.1"
        },
        {
            "bk_property_id": "bk_host_outerip",
            "bk_property_name": "外网IP",
            "bk_property_value": ""
        },
		......
        {
            "bk_property_id": "bk_addressing",
            "bk_property_name": "寻址方式",
            "bk_property_value": "static"
        }
    ]
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
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| bk_property_id    | string     | 属性id |
| bk_property_name  | string     | 属性名称 |
| bk_property_value | string     | 属性值 |
